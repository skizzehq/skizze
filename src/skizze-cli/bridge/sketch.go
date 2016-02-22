package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	pb "datamodel/protobuf"

	"github.com/gogo/protobuf/proto"
)

func createSketch(fields []string, in *pb.Sketch) error {
	if in.GetType() != pb.SketchType_CARD {
		if len(fields) > 4 {
			return fmt.Errorf("Too many argumets, expected 4 got %d", len(fields))
		}
		num, err := strconv.Atoi(fields[3])
		if err != nil {
			return fmt.Errorf("Expected last argument to be of type int: %q", err)
		}

		in.Properties = &pb.SketchProperties{
			Size:           proto.Int64(int64(num)),
			MaxUniqueItems: proto.Int64(int64(num)),
		}
	}
	_, err := client.CreateSketch(context.Background(), in)
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

	switch strings.ToLower(fields[0]) {
	case "create":
		return createSketch(fields, in)
	case "add":
		return addToSketch(fields, in)
	case "get":
		return getFromSketch(fields, in)
	case "destroy":
	case "info":
		return getSketchInfo(in)
	default:
		return fmt.Errorf("unkown operation: %s", fields[0])
	}
	return nil
}

func listSketches() error {
	reply, err := client.ListAll(context.Background(), &pb.Empty{})
	if err == nil {
		for _, v := range reply.GetSketches() {
			line := fmt.Sprintf("Name: %s\t  Type: %s", v.GetName(), v.GetType().String())
			_, _ = fmt.Fprintln(w, line)
		}
		_ = w.Flush()
	}
	return err
}

func listSketchType(typ pb.SketchType) error {
	reply, err := client.List(context.Background(), &pb.ListRequest{Type: &typ})
	if err == nil {
		for _, v := range reply.GetSketches() {
			line := fmt.Sprintf("Name: %s\t  Type: %s", v.GetName(), v.GetType().String())
			_, _ = fmt.Fprintln(w, line)
		}
		_ = w.Flush()
	}
	return err
}

func getSketchInfo(in *pb.Sketch) error {
	in.Properties = &pb.SketchProperties{
		MaxUniqueItems: proto.Int64(0),
		Size:           proto.Int64(0),
	}
	reply, err := client.GetSketch(context.Background(), in)
	logger.Infof("%s", reply)
	return err
}

func getFromSketch(fields []string, in *pb.Sketch) error {
	if len(fields) < 3 {
		return fmt.Errorf("Expected at least 3 values, got %q", len(fields))
	}
	getRequest := &pb.GetRequest{
		Sketches: []*pb.Sketch{in},
		Values:   fields[3:],
	}

	switch in.GetType() {
	case pb.SketchType_CARD:
		reply, err := client.GetCardinality(context.Background(), getRequest)
		if err == nil {
			if len(reply.GetResults()) == 0 {
				logger.Infof("%s does not exist", in.GetName())
			} else {
				logger.Infof("Cardinality: %s", reply.GetResults()[0].GetCardinality())
			}
		}
		return err
	case pb.SketchType_FREQ:
		reply, err := client.GetFrequency(context.Background(), getRequest)
		if err == nil {
			if len(reply.GetResults()) == 0 {
				logger.Infof("%s does not exist", in.GetName())
			} else {
				for _, v := range reply.GetResults()[0].GetFrequencies() {
					line := fmt.Sprintf("Value: %s\t  Hits: %d", v.GetValue(), v.GetCount())
					_, _ = fmt.Fprintln(w, line)
				}
			}
			_ = w.Flush()
		}
		return err
	case pb.SketchType_MEMB:
		reply, err := client.GetMembership(context.Background(), getRequest)
		if err == nil {
			if len(reply.GetResults()) == 0 {
				logger.Infof("%s does not exist", in.GetName())
			} else {
				for _, v := range reply.GetResults()[0].GetMemberships() {
					line := fmt.Sprintf("Value: %s\t  Member: %t", v.GetValue(), v.GetIsMember())
					_, _ = fmt.Fprintln(w, line)
				}
				_ = w.Flush()
			}
		}
		return err
	case pb.SketchType_RANK:
		reply, err := client.GetRankings(context.Background(), getRequest)
		if err == nil {
			if len(reply.GetResults()) == 0 {
				logger.Infof("%s does not exist", in.GetName())
			} else {
				for i, v := range reply.GetResults()[0].GetRankings() {
					line := fmt.Sprintf("Rank: %d\t  Value: %s\t  Hits: %d", i+1, v.GetValue(), v.GetCount())
					_, _ = fmt.Fprintln(w, line)
				}
				_ = w.Flush()
			}
		}
		return err
	default:
		return fmt.Errorf("Unkown Type %s", in.GetType().String())
	}
}
