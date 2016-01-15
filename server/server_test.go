package server

import (
	"log"
	"time"

	pb "github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/manager"
	"google.golang.org/grpc"
)

func setupClient() (pb.SkizzeClient, *grpc.ClientConn) {
	m := manager.NewManager()
	go Run(m, 7777)
	time.Sleep(time.Millisecond * 50)

	// Connect to the server.
	conn, err := grpc.Dial("127.0.0.1:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return pb.NewSkizzeClient(conn), conn
}

func tearDownClient(conn *grpc.ClientConn) {
	conn.Close()
	Stop()
}
