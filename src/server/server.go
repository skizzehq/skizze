package server

import (
	"config"
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	pb "datamodel/protobuf"
	"manager"
	"storage"
	"utils"

	"golang.org/x/net/context"
)

type serverStruct struct {
	manager *manager.Manager
	g       *grpc.Server
	storage *storage.AOF
}

var server *serverStruct

// Run ...
func Run(manager *manager.Manager, port uint) {
	path := filepath.Join(config.GetConfig().DataDir, "skizze.aof")
	storage := storage.NewAOF(path)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // RPC port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g := grpc.NewServer()

	server = &serverStruct{manager, g, storage}
	pb.RegisterSkizzeServer(g, server)
	server.replay()
	_ = g.Serve(lis)
}

func unmarshalSketch(e *storage.Entry) *pb.Sketch {
	sketch := &pb.Sketch{}
	err := proto.Unmarshal(e.Args(), sketch)
	utils.PanicOnError(err)
	return sketch
}

func unmarshalDom(e *storage.Entry) *pb.Domain {
	dom := &pb.Domain{}
	err := proto.Unmarshal(e.Args(), dom)
	utils.PanicOnError(err)
	return dom
}

func (server *serverStruct) replay() {
	fmt.Println("Replaying ...")
	for {
		e, err := server.storage.Read()
		if err != nil && err.Error() == "EOF" {
			break
		} else {
			utils.PanicOnError(err)
		}

		switch e.Op() {
		case storage.Add:
			req := &pb.AddRequest{}
			err = proto.Unmarshal(e.Args(), req)
			utils.PanicOnError(err)
			if _, err := server.add(context.Background(), req); err != nil {
				fmt.Println(err)
			}
		case storage.CreateSketch:
			sketch := unmarshalSketch(e)
			if _, err := server.createSketch(context.Background(), sketch); err != nil {
				fmt.Println(err)
			}
		case storage.DeleteSketch:
			sketch := unmarshalSketch(e)
			if _, err := server.deleteSketch(context.Background(), sketch); err != nil {
				fmt.Println(err)
			}
		case storage.CreateDom:
			dom := unmarshalDom(e)
			if _, err := server.createDomain(context.Background(), dom); err != nil {
				fmt.Println(err)
			}
		case storage.DeleteDom:
			dom := unmarshalDom(e)
			if _, err := server.deleteDomain(context.Background(), dom); err != nil {
				fmt.Println(err)
			}
		default:
			continue
		}
	}
}

// Stop ...
func Stop() {
	server.g.Stop()
}
