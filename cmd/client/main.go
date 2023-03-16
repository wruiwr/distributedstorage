package main

import (
	"bytes"
	"flag"
	"fmt"
	// qc "github.com/selabhvl/cpnmbt/rwregister"
	qc "github.com/wruiwr/distributedstorage"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	qSize = 2 // assumed quorum size is 2
	n     = 3 // the number of local listening servers

)

func main() {

	var (
		port   = flag.Int("port", 8080, "port where local server is listening")
		saddrs = flag.String("addrs", "", "server addresses separated by ','")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *saddrs == "" {

		var buf bytes.Buffer
		for i := 0; i < n; i++ {
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(*port + i))
			buf.WriteString(",")
		}
		b := buf.String()
		*saddrs = b[:len(b)-1]
	}

	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}
	log.Printf("the number of local listening servers: %d, addrs:(%v)", len(addrs), *saddrs)

	//options:
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		grpc.WithInsecure(),
	}
	dialOpts := qc.WithGrpcDialOptions(grpcOpts...)

	mgrOpts := []qc.ManagerOption{
		dialOpts,
		qc.WithTracing(),
	}

	mgr, err := qc.NewManager(
		addrs,
		mgrOpts...,
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer mgr.Close()

	// Get all all available node ids, 3 nodes for example:
	ids := mgr.NodeIDs()

	fmt.Println("Nodes ids:", ids)

	// Quorum Specs
	qspec := qc.NewQuoSpecQ(qSize, qSize)

	// The configuration for both Write and Read,
	// using nodes ids and the instance of NewQuoSpecQ (QuorumSpec) as the input
	config, err := mgr.NewConfiguration(ids, qspec)
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	// A test state
	state := &qc.Value{
		C: &qc.Content{Value: "Rui", Timestamp: time.Now().UnixNano()},
	}
	// Perform write call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	wreply, err := config.Write(ctx, state)
	if err != nil {
		log.Fatalf("write quorum call error: %v", err)
	}

	fmt.Printf("wreply: %v\n", wreply)
	if !wreply.Ack {
		//t.Error("write reply was not marked as new")
		fmt.Println("write reply was not marked as ack")
	}

	// Perform read call
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rreply, err := config.Read(ctx, &qc.ReadRequest{})
	if err != nil {
		log.Fatalf("read quorum call error: %v", err)
	}
	fmt.Printf("rreply: %v\n", rreply) // to see the result
	fmt.Printf("rreply: %v\n", rreply.C.Value)
	/*	if rreply.Value.C.Value != state.C.Value {
		log.Fatalf("read reply: got value %v, want value %v", rreply.C.Value, state.C.Value)
	}*/

	nodes := mgr.Nodes()
	for _, m := range nodes {
		fmt.Printf("%v\n", m)
	}
}
