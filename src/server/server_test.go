package server

import (
	"time"

	"google.golang.org/grpc"

	pb "datamodel/protobuf"
	"manager"
)

func setupClient() (pb.SkizzeClient, *grpc.ClientConn) {
	m := manager.NewManager()
	go Run(m, "127.0.0.1", 7777)
	time.Sleep(time.Millisecond * 50)

	// Connect to the server.
	conn, err := grpc.Dial("127.0.0.1:7777", grpc.WithInsecure())
	if err != nil {
		logger.Criticalf("fail to dial: %v", err)
	}
	return pb.NewSkizzeClient(conn), conn
}

func tearDownClient(conn *grpc.ClientConn) {
	_ = conn.Close()
	Stop()
}
