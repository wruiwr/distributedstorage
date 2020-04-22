package rwregister

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"strconv"

	r "github.com/selabhvl/cpnmbt/rwregister/reader"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const allServers = -1

var (
	sysTests     r.Test
	portSupplier = struct {
		p int
		sync.Mutex
	}{}
	serviceTerminated = false
)

func TestMain(m *testing.M) {
	// Flag definitions.
	var hosts = flag.String(
		"remotehosts",
		"",
		"comma separated list of 'addr:port' pairs to use as hosts for remote benchmarks",
	)
	var portBase = flag.Int(
		"portbase",
		22332,
		"use a specific port base (incremented for each needed test listener)",
	)
	var dir = flag.String(
		"dir",
		"tests.xml", // this is the latest xml format
		"path to system test file",
	)

	// Parse and validate flags.
	flag.Parse()
	err := parseHostnames(*hosts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	portSupplier.p = *portBase

	// Load the system test cases from XML file
	r.ParseXMLTestCase(*dir, &sysTests)

	// Enable gRPC tracing and logging.
	silentLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	grpclog.SetLogger(silentLogger)
	grpc.EnableTracing = true

	// Run tests/benchmarks.
	res := m.Run()
	os.Exit(res)
}

func TestNoConnection(t *testing.T) {
	// Tests that the creating manager fails when connections are down
	var storServers storageServers
	for i := 0; i < 3; i++ {
		storServers = append(storServers, storageServer{impl: NewRegisterBasic()})
	}
	servers, dialOpts, stopGrpcServe, closeListeners := setup(t, storServers, false, false)
	closeListeners(allServers)
	stopGrpcServe(allServers)
	mgr, err := NewManager(servers.addrs(), dialOpts, WithTracing())
	if err == nil || mgr != nil {
		t.Fatalf("expected create manager to fail")
	}
}

func TestSystemQuorumCalls(t *testing.T) {

	t.Logf("%v (system size=%d, quorum size=%d)", sysTests.Name, sysTests.SystemSize, sysTests.QuorumSize)

	//TODO: Need to consider server failures.

	for _, testcase := range sysTests.TestCases {
		t.Logf("%s, Description: %s", testcase.Name, testcase.Description)

		mgr, config, stopGrpcServe, closeListeners := setupConfig(t, sysTests.SystemSize, sysTests.QuorumSize)
		defer mgr.Close()
		defer closeListeners(allServers)
		defer stopGrpcServe(allServers)

		eventMonitor := NewEventMonitor()
		operation := &operations{t, eventMonitor, config, closeListeners, stopGrpcServe}

		for _, ops := range testcase.OrderOp {
			switch ops.OpType {
			case "Concurrent":
				// concurrent execution of the test cases
				var wg sync.WaitGroup
				wg.Add(len(ops.Op))
				for _, routine := range ops.Op {
					go func(routine r.Operation) {
						funcExecutor(t, operation, routine.Name, routine.ID, routine.Value)
						wg.Done()
					}(routine)
				}
				// wait for all concurrent goroutines to finish before continuing test.
				wg.Wait()

			case "Sequential":
				// sequential execution of the test cases
				for _, routine := range ops.Op {
					funcExecutor(t, operation, routine.Name, routine.ID, routine.Value)
				}
			}
		}
	}
}

// funcExecutor is used to executed different functions.
func funcExecutor(t *testing.T, oper interface{}, funcName string, params ...interface{}) {
	t.Helper()
	inputArgs := make([]reflect.Value, len(params))
	for i, param := range params {
		inputArgs[i] = reflect.ValueOf(param)
	}
	fn := reflect.ValueOf(oper).MethodByName(funcName)
	if !fn.IsValid() {
		t.Errorf("method '%s' not found", funcName)
	}
	fn.Call(inputArgs)
}

// operations represent the executable operations of the tests.
type operations struct {
	t              *testing.T
	monitor        Monitor
	config         *Configuration
	closeListeners func(n int)
	stopGrpcServe  func(n int)
}

// createValue can create a value for Write quorum call from the write value of xml file.
func createValue(writeValue string) *Value {
	return &Value{
		C: &Content{Value: writeValue, Timestamp: time.Now().UnixNano()},
	}
}

// serverFailure terminates k servers. If more than a quorum of servers are terminated then
// global variable 'serviceTerminated' will be set to true to indicate that future quorum
// calls will not return a value, but an incomplete error instead.
func (op operations) ServerFailure(id, params string) {
	k, _ := strconv.Atoi(params)
	for n := 1; n <= k; n++ {
		op.closeListeners(n)
		op.stopGrpcServe(n)
		op.t.Logf("The server %d is terminated", n)
	}
	if k >= sysTests.QuorumSize {
		serviceTerminated = true
	}
}

// DoWriteCall preforms a Write quorum call and tests the value to be stored against
// the set of legal output values obtained from the test oracle/monitor.
func (op operations) DoWriteCall(routineID string, writeValue string) {
	value := createValue(writeValue)
	if value == nil {
		op.t.Fatal("Cannot invoke Write with nil value")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	i := op.monitor.WI(value.C.Value)
	wreply, err := op.config.Write(ctx, value)
	op.monitor.WR(i, value.C.Value)

	if serviceTerminated {
		// a quorum of servers have been terminated, expect a quorum call error
		if err == nil {
			if wreply != nil {
				op.t.Errorf("expected quorum call error, but got %v", wreply)
			} else {
				op.t.Errorf("expected quorum call error")
			}
		}
		return
	}
	// we have a quorum of servers running, expect Ack = true, without any errors
	if err != nil {
		op.t.Errorf("Routine: %v, Write quorum call error: %v", routineID, err)
	}
	// no need for nil check on wreply; GetAck() does the check.
	if !wreply.GetAck() {
		op.t.Errorf("Got Write ACK: %v, want %v", wreply.GetAck(), !wreply.GetAck())
	}
}

// DoReadCall preforms a Read quorum call and tests the reply value against
// the set of legal output values obtained from the test oracle/monitor.
// For the failure scenarios, if we inject an error in quorum functions,
// such as "len(replies) <= qq.wq", after terminating one of servers,
// the test adaptor can capture the error.
// For the DoReadCall function, readValue argument is not used for the read call.
func (op operations) DoReadCall(routineID string, readValue string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	i := op.monitor.RI()
	rreply, err := op.config.Read(ctx, &ReadRequest{})

	// check if the server has already been terminated.
	if !serviceTerminated {
		// before servers are terminated
		if err != nil {
			op.t.Errorf("Routine: %v, Read quorum call error: %v", routineID, err)
		}
		if rreply == nil {
			op.t.Errorf("Read reply: %v, got an error: %v", rreply, err)
		}
		v := rreply.C.Value
		isOK, readLegalValues := op.monitor.RR(i, v)
		if !isOK {
			op.t.Errorf("Read quorum call reply: got %v, not within the wanted legal value list %v", v, readLegalValues)
		} else {
			op.t.Logf("Read quorum call reply: got %v, within the wanted legal value list %v \n", v, readLegalValues)
		}
	} else {
		// after servers are terminated
		if err == nil && rreply != nil {
			op.t.Errorf("expect quorum call error: incomplete call")
		}
	}
}

func setupConfig(t testing.TB, systemSize, quorumSize int) (*Manager, *Configuration, func(n int), func(n int)) {
	t.Helper()
	var storServers storageServers
	for i := 0; i < systemSize; i++ {
		storServers = append(storServers, storageServer{impl: NewRegisterBasic()})
	}
	servers, dialOpts, stopGrpcServe, closeListeners := setup(t, storServers, false, false)

	mgr, err := NewManager(servers.addrs(), dialOpts, WithTracing())
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	ids := mgr.NodeIDs()
	t.Logf("Nodes ids: %v", ids)
	config, err := mgr.NewConfiguration(ids, NewQuoSpecQ(quorumSize, quorumSize))
	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}
	return mgr, config, stopGrpcServe, closeListeners
}

func setup(t testing.TB, storServers []storageServer, remote, secure bool) (storageServers, ManagerOption, func(n int), func(n int)) {
	t.Helper()
	if len(storServers) == 0 {
		t.Fatal("setupServers: need at least one server")
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
	}
	grpcOpts = append(grpcOpts, grpc.WithInsecure())
	dialOpts := WithGrpcDialOptions(grpcOpts...)

	if remote {
		return storServers, dialOpts, func(int) {}, func(int) {}
	}

	servers := make([]*grpc.Server, len(storServers))
	for i := range servers {
		var opts []grpc.ServerOption
		servers[i] = grpc.NewServer(opts...)
		RegisterRegisterServer(servers[i], storServers[i].impl)
		if storServers[i].addr == "" {
			portSupplier.Lock()
			storServers[i].addr = fmt.Sprintf("localhost:%d", portSupplier.p)
			portSupplier.p++
			portSupplier.Unlock()
		}
	}

	listeners := make([]net.Listener, len(servers))

	var err error
	for i, rs := range storServers {
		listeners[i], err = net.Listen("tcp", rs.addr)
		if err != nil {
			t.Fatalf("failed to listen: %v", err)
		}
	}
	for i, server := range servers {
		go func(i int, server *grpc.Server) {
			_ = server.Serve(listeners[i])
		}(i, server)
	}

	stopGrpcServeFunc := func(n int) {
		if n < 0 || n > len(servers) {
			for _, s := range servers {
				s.Stop()
			}
		} else {
			servers[n].Stop()
		}
	}

	closeListenersFunc := func(n int) {
		if n < 0 || n > len(listeners) {
			for _, l := range listeners {
				l.Close()
			}
		} else {
			listeners[n].Close()
		}
	}

	return storServers, dialOpts, stopGrpcServeFunc, closeListenersFunc
}

type storageServer struct {
	impl RegisterTestServer
	addr string
}

type storageServers []storageServer

func (rs storageServers) addrs() []string {
	addrs := make([]string, len(rs))
	for i, server := range rs {
		addrs[i] = server.addr
	}
	return addrs
}

var remoteBenchmarkHosts []string

func parseHostnames(hostnames string) error {
	if hostnames == "" {
		return nil
	}

	hostPairsSplitted := strings.Split(hostnames, ",")
	for i, hps := range hostPairsSplitted {
		tmp := strings.Split(hps, ":")
		if len(tmp) != 2 {
			return fmt.Errorf("parseHostnames: malformed host address: host %d: %q", i, hps)
		}
		remoteBenchmarkHosts = append(remoteBenchmarkHosts, hps)
	}

	return nil
}
