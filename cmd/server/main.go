package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	// qc "github.com/selabhvl/cpnmbt/rwregister"
	qc "github.com/wruiwr/distributedstorage"
	"google.golang.org/grpc"
)

func main() {
	var port = flag.Int("port", 8080, "port to listen on")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	server(*port)
}

type regServer struct {
	impl qc.RegisterTestServer
	addr string
}

func server(port int) {
	var rs = regServer{
		impl: qc.NewRegisterBasic(),
		addr: fmt.Sprintf("localhost:%d", port),
	}
	l, err := net.Listen("tcp", rs.addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	qc.RegisterRegisterServer(grpcServer, rs.impl)

	log.Printf("server %s running", rs.addr)
	log.Fatal(grpcServer.Serve(l))
}
