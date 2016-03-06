package server

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"

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
func Run(manager *manager.Manager, host string, port int, datadir string) {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	path := filepath.Join(datadir, "skizze.aof")
	aof := storage.NewAOF(path)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port)) // RPC port
	if err != nil {
		logger.Criticalf("failed to listen: %v", err)
	}
	g := grpc.NewServer()

	server = &serverStruct{manager, g, aof}
	pb.RegisterSkizzeServer(g, server)
	server.replay()
	aof.Run()
	_ = g.Serve(lis)
}

func unmarshalSketch(e *storage.Entry) *pb.Sketch {
	sketch := &pb.Sketch{}
	err := proto.Unmarshal(e.RawMsg(), sketch)
	utils.PanicOnError(err)
	return sketch
}

func unmarshalDom(e *storage.Entry) *pb.Domain {
	dom := &pb.Domain{}
	err := proto.Unmarshal(e.RawMsg(), dom)
	utils.PanicOnError(err)
	return dom
}

func (server *serverStruct) replay() {
	logger.Infof("Replaying ...")
	for {
		e, err := server.storage.Read()
		if err != nil && err.Error() == "EOF" {
			break
		} else {
			utils.PanicOnError(err)
		}

		switch e.OpType() {
		case storage.Add:
			req := &pb.AddRequest{}
			err = proto.Unmarshal(e.RawMsg(), req)
			utils.PanicOnError(err)
			if _, err := server.add(context.Background(), req); err != nil {
				logger.Errorf("an error has occurred while replaying: %s", err.Error())
			}
		case storage.CreateSketch:
			sketch := unmarshalSketch(e)
			if _, err := server.createSketch(context.Background(), sketch); err != nil {
				logger.Errorf("an error has occurred while replaying: %s", err.Error())
			}
		case storage.DeleteSketch:
			sketch := unmarshalSketch(e)
			if _, err := server.deleteSketch(context.Background(), sketch); err != nil {
				logger.Errorf("an error has occurred while replaying: %s", err.Error())
			}
		case storage.CreateDom:
			dom := unmarshalDom(e)
			if _, err := server.createDomain(context.Background(), dom); err != nil {
				logger.Errorf("an error has occurred while replaying: %s", err.Error())
			}
		case storage.DeleteDom:
			dom := unmarshalDom(e)
			if _, err := server.deleteDomain(context.Background(), dom); err != nil {
				logger.Errorf("an error has occurred while replaying: %s", err.Error())
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
