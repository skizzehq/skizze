package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	pb "datamodel/protobuf"

	"github.com/gogo/protobuf/proto"
)

func createDomain(fields []string, in *pb.Domain) error {
	if len(fields) != 5 {
		return fmt.Errorf("Expected 5 arguments got %d", len(fields))
	}

	// FIXME make last 2 arguments optional
	capa, err := strconv.Atoi(fields[3])
	if err != nil {
		return fmt.Errorf("Expected 3rd argument to be of type int: %q", err)
	}

	size, err := strconv.Atoi(fields[4])
	if err != nil {
		return fmt.Errorf("Expected last argument to be of type int: %q", err)
	}

	types := []pb.SketchType{pb.SketchType_MEMB, pb.SketchType_FREQ, pb.SketchType_RANK, pb.SketchType_CARD}
	for _, ty := range types {
		sketch := &pb.Sketch{}
		sketch.Name = proto.String("")
		sketch.Type = &ty
		sketch.Properties = &pb.SketchProperties{
			Size:           proto.Int64(int64(size)),
			MaxUniqueItems: proto.Int64(int64(capa)),
		}
		in.Sketches = append(in.Sketches, sketch)
	}

	_, err = client.CreateDomain(context.Background(), in)
	if err == nil {
		fmt.Println("done")
	}
	return err
}

func addToDomain(fields []string, in *pb.Domain) error {
	if len(fields) < 4 {
		return fmt.Errorf("Expected at least 4 values, got %q", len(fields))
	}
	addRequest := &pb.AddRequest{
		Domain: in,
		Values: fields[3:],
	}
	_, err := client.Add(context.Background(), addRequest)
	if err == nil {
		fmt.Println("done")
	}
	return err
}

func sendDomainRequest(fields []string) error {
	name := fields[2]
	in := &pb.Domain{
		Name: proto.String(name),
	}

	switch strings.ToLower(fields[0]) {
	case "create":
		return createDomain(fields, in)
	case "add":
		return addToDomain(fields, in)
	case "destroy":
		return deleteDomain(fields, in)
	case "info":
		return getDomainInfo(fields, in)
	default:
		return fmt.Errorf("unkown operation: %s", fields[0])
	}
}

func listDomains() error {
	reply, err := client.ListDomains(context.Background(), &pb.Empty{})
	if err == nil {
		for _, v := range reply.GetNames() {
			_, _ = fmt.Fprintln(w, fmt.Sprintf("Name: %s\t", v))
		}
		_ = w.Flush()
	}
	return err
}

func deleteDomain(fields []string, in *pb.Domain) error {
	_, err := client.DeleteDomain(context.Background(), in)
	return err
}

func getDomainInfo(fields []string, in *pb.Domain) error {
	dom, err := client.GetDomain(context.Background(), in)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(w, fmt.Sprintf("Name: %s  Type: %s\t", dom.GetName(), ""))
	_, _ = fmt.Fprintln(w, fmt.Sprintf("%d Sketches attached:", len(dom.GetSketches())))
	for i, v := range dom.GetSketches() {
		_, _ = fmt.Fprintln(w, fmt.Sprintf("  %d.  Name: %s  Type: %s\t", i+1, v.GetName(), v.GetType()))
	}
	_ = w.Flush()
	return err
}
