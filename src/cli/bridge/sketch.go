package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/gogo/protobuf/proto"
	pb "github.com/skizzehq/skizze/datamodel"
)

func createSketch(fields []string, in *pb.Sketch) error {
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

func addToSketch(fields []string, in *pb.Sketch) error {
	if len(fields) < 4 {
		return fmt.Errorf("Expected at least 4 values, got %q", len(fields))
	}
	addRequest := &pb.AddRequest{
		Sketch: in,
		Values: fields[3:],
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

	switch strings.ToLower(fields[1]) {
	case "create":
		return createSketch(fields, in)
	case "add":
		return addToSketch(fields, in)
	case "get":
		return getFromSketch(fields, in)
	case "destroy":
	case "info":
	default:
		return fmt.Errorf("unkown operation: %s", fields[1])
	}
	return nil
}

func getFromSketch(fields []string, in *pb.Sketch) error {
	if len(fields) < 3 {
		return fmt.Errorf("Expected at least 3 values, got %q", len(fields))
	}
	getRequest := &pb.GetRequest{
		Sketch: in,
		Values: fields[3:],
	}

	switch in.GetType() {
	case pb.SketchType_CARD:
		reply, err := client.GetCardinality(context.Background(), getRequest)
		if err == nil {
			fmt.Println(reply)
		}
		return err
	case pb.SketchType_FREQ:
		reply, err := client.GetFrequency(context.Background(), getRequest)
		if err == nil {
			fmt.Println(reply)
		}
		return err
	case pb.SketchType_MEMB:
		reply, err := client.GetMembership(context.Background(), getRequest)
		if err == nil {
			fmt.Println(reply)
		}
		return err
	case pb.SketchType_RANK:
		reply, err := client.GetRank(context.Background(), getRequest)
		if err == nil {
			fmt.Println(reply)
		}
		return err
	default:
		return fmt.Errorf("Unkown Type %s", in.GetType().String())
	}
}
