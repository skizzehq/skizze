package bridge

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	pb "datamodel"

	"github.com/peterh/liner"
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
			fields = append(fields, f)
		}
	}
	return fields
}

func evalutateQuery(query string) error {
	fields := getFields(query)
	switch strings.ToLower(fields[0]) {
	case pb.HLLPP:
		return sendSketchRequest(fields, pb.SketchType_CARD)
	case pb.CML:
		return sendSketchRequest(fields, pb.SketchType_FREQ)
	case pb.TopK:
		return sendSketchRequest(fields, pb.SketchType_RANK)
	case pb.Bloom:
		return sendSketchRequest(fields, pb.SketchType_MEMB)
	case pb.DOM:
		return sendDomainRequest(fields)
	default:
		return fmt.Errorf("unkown field or command %s", fields[0])
	}
}

// Run ...
func Run() {
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
