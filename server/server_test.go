package server

import (
	"log"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/gogo/protobuf/proto"
	"github.com/seiflotfy/skizze/config"
	pb "github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/manager"
	"github.com/seiflotfy/skizze/utils"
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

func TestCreateSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}
}

func TestGetAddCardSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d"},
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	getReq := &pb.GetRequest{
		Sketch: in,
		Values: []string{},
	}

	if res, err := client.GetCardinality(context.Background(), getReq); err != nil {
		t.Error("Did not expect error, got", err)

	} else if res.GetCardinality() != 4 {
		t.Error("Expected cardinality 4, got", res.GetCardinality())
	}
}
