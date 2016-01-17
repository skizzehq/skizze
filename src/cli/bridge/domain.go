package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/gogo/protobuf/proto"
	pb "github.com/skizzehq/skizze/datamodel"
)

func createDomain(fields []string, in *pb.Domain) error {
	if len(fields) > 5 {
		return fmt.Errorf("Too many argumets, expected 4 got %d", len(fields))
	}

	// FIXME make last 2 arguments optional
	capa, err := strconv.Atoi(fields[3])
	if err != nil {
		return fmt.Errorf("Expected 3rd argument to be of type int: %q", err)
	}

	rank, err := strconv.Atoi(fields[4])
	if err != nil {
		return fmt.Errorf("Expected last argument to be of type int: %q", err)
	}

	in.Defaults = &pb.Defaults{
		Rank:     proto.Int64(int64(rank)),
		Capacity: proto.Int64(int64(capa)),
	}

	_, err = client.CreateDomain(context.Background(), in)
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
	return err
}

func sendDomainRequest(fields []string) error {
	name := fields[2]
	typ := pb.DomainType_STRING
	in := &pb.Domain{
		Name: proto.String(name),
		Type: &typ,
	}

	switch strings.ToLower(fields[1]) {
	case "create":
		return createDomain(fields, in)
	case "add":
		return addToDomain(fields, in)
	//case "destroy":
	//case "info":
	default:
		return fmt.Errorf("unkown operation: %s", fields[1])
	}
}
