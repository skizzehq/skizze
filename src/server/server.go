package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "datamodel"
	"manager"
)

type serverStruct struct {
	manager *manager.Manager
	g       *grpc.Server
}

var server *serverStruct

// Run ...
func Run(manager *manager.Manager, port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // RPC port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g := grpc.NewServer()

	server = &serverStruct{manager, g}

	pb.RegisterSkizzeServer(g, server)
	_ = g.Serve(lis)
}

// Stop ...
func Stop() {
	server.g.Stop()
}
