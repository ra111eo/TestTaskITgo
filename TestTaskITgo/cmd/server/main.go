package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ra111eo/TestTaskITgo/pb"
	"github.com/ra111eo/TestTaskITgo/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	ewalletServer := service.NewEwalletServer(service.NewDBEwalletStore())
	log.Println("store created succsessfully")
	grcpServer := grpc.NewServer()
	pb.RegisterEwalletServiceServer(grcpServer, ewalletServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("can not start server:", err)
	}

	err = grcpServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
