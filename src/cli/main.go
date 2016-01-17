package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/gogo/protobuf/proto"
	"github.com/peterh/liner"

	pb "datamodel"
)

var client pb.SkizzeClient
var conn *grpc.ClientConn
var historyFn = filepath.Join(os.TempDir(), ".skizze_history")

func setupClient() (pb.SkizzeClient, *grpc.ClientConn) {
	// Connect to the server.
	conn, err := grpc.Dial("127.0.0.1:3596", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return pb.NewSkizzeClient(conn), conn
}

func tearDownClient(conn *grpc.ClientConn) {
	conn.Close()
}

func getFields(query string) []string {
	fields := []string{}
	for _, f := range strings.Split(query, " ") {
		if len(f) > 0 {
			fields = append(fields, strings.ToLower(f))
		}
	}
	return fields
}

func evalutateQuery(query string) error {
	fields := getFields(query)
	switch fields[0] {
	case pb.HLLPP:
		return sendSketchRequest(fields, pb.SketchType_CARD)
	case pb.CML:
		return sendSketchRequest(fields, pb.SketchType_FREQ)
	case pb.TopK:
		return sendSketchRequest(fields, pb.SketchType_RANK)
	case pb.Bloom:
		return sendSketchRequest(fields, pb.SketchType_MEMB)
	default:
		return fmt.Errorf("unkown field or command %s", fields[0])
	}
}

func sendCreateSketchRequest(fields []string, in *pb.Sketch) error {
	if len(fields) > 4 {
		return fmt.Errorf("Too many argumets, expected 4 got %d", len(fields))
	}

	num, err := strconv.Atoi(fields[3])
	if err != nil {
		return fmt.Errorf("Expected last argument to be of type int: %q", err)
	}

	in.Defaults = &pb.Defaults{
		Rank:     proto.Int64(int64(num)),
		Capacity: proto.Int64(int64(num)),
	}

	_, err = client.CreateSketch(context.Background(), in)
	return err
}

func sendAddSketchRequest(fields []string, in *pb.Sketch) error {
	if len(fields) < 4 {
		return fmt.Errorf("Expected at least 4 values, got %q", len(fields))
	}
	addRequest := &pb.AddRequest{
		Sketch: in,
		Values: fields[2:],
	}
	_, err := client.Add(context.Background(), addRequest)
	return err
}

func sendSketchRequest(fields []string, typ pb.SketchType) error {
	name := fields[2]
	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	switch fields[1] {
	case "create":
		return sendCreateSketchRequest(fields, in)
	case "info":
	case "destroy":
	case "add":
		return sendAddSketchRequest(fields, in)
	default:
		return fmt.Errorf("unkown operation: %s", fields[1])
	}
	return nil
}

func main() {
	client, conn = setupClient()
	line := liner.NewLiner()
	defer func() { _ = line.Close() }()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFn); err == nil {
		if _, err := line.ReadHistory(f); err == nil {
			_ = f.Close()
		}
	}

	for {
		if query, err := line.Prompt("skizze-cli> "); err == nil {
			if err := evalutateQuery(query); err != nil {
				fmt.Println(err)
			}
			line.AppendHistory(query)
		} else if err == liner.ErrPromptAborted {
			log.Print("Aborted")
			tearDownClient(conn)
			return
		} else {
			log.Print("Error reading line: ", err)
		}

		if f, err := os.Create(historyFn); err != nil {
			log.Print("Error writing history file: ", err)
		} else {
			if _, err := line.WriteHistory(f); err != nil {
				_ = f.Close()
			}
		}
	}
}
