package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/manager"
	"google.golang.org/grpc"
)

type server struct {
	manager *manager.Manager
}

// Run ...
func Run(manager *manager.Manager, port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // RPC port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g := grpc.NewServer()

	server := &server{manager}

	pb.RegisterSkizzeServer(g, server)
	g.Serve(lis)
}
