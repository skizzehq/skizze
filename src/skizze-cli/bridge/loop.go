package bridge

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"datamodel"
	pb "datamodel/protobuf"

	"github.com/peterh/liner"
)

var client pb.SkizzeClient
var conn *grpc.ClientConn
var historyFn = filepath.Join(os.TempDir(), ".skizze_history")
var w = new(tabwriter.Writer)
var typeMap = map[string]pb.SketchType{
	datamodel.HLLPP: pb.SketchType_CARD,
	datamodel.CML:   pb.SketchType_FREQ,
	datamodel.Bloom: pb.SketchType_MEMB,
	datamodel.TopK:  pb.SketchType_RANK,
}

func setupClient() (pb.SkizzeClient, *grpc.ClientConn) {
	// Connect to the server.
	var err error
	conn, err = grpc.Dial("127.0.0.1:3596", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return pb.NewSkizzeClient(conn), conn
}

func tearDownClient(conn io.Closer) {
	_ = conn.Close()
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
	if len(fields) <= 2 {
		//TODO: global stuff might be set
		switch strings.ToLower(fields[0]) {
		case "list":
			if len(fields) == 1 {
				return listSketches()
			} else if len(fields) == 2 && strings.ToLower(fields[1]) == "dom" {
				return listDomains()
			} else if len(fields) == 2 {
				v, ok := typeMap[strings.ToLower(fields[1])]
				if !ok {
					return fmt.Errorf("Invalid operation: %s", query)
				}
				return listSketchType(v)
			}
		case "save":
			if len(fields) == 1 {
				return save()
			}
			return fmt.Errorf("Invalid operation: %s", query)
		default:
			return fmt.Errorf("Invalid operation: %s", query)
		}
	}

	if len(fields) > 2 {
		switch strings.ToLower(fields[1]) {
		case datamodel.HLLPP:
			return sendSketchRequest(fields, pb.SketchType_CARD)
		case datamodel.CML:
			return sendSketchRequest(fields, pb.SketchType_FREQ)
		case datamodel.TopK:
			return sendSketchRequest(fields, pb.SketchType_RANK)
		case datamodel.Bloom:
			return sendSketchRequest(fields, pb.SketchType_MEMB)
		case datamodel.DOM:
			return sendDomainRequest(fields)
		default:
			return fmt.Errorf("unkown field or command %s", fields[1])
		}
	}
	return errors.New("Invalid operation")
}

func save() error {
	_, err := client.CreateSnapshot(context.Background(), &pb.CreateSnapshotRequest{})
	return err
}

// Run ...
func Run() {
	client, conn = setupClient()
	line := liner.NewLiner()
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	defer func() { _ = line.Close() }()

	line.SetCtrlCAborts(true)

	if f, err := os.Open(historyFn); err == nil {
		if _, err := line.ReadHistory(f); err == nil {
			_ = f.Close()
		}
	}

	for {
		if query, err := line.Prompt("skizze> "); err == nil {
			if err := evalutateQuery(query); err != nil {
				fmt.Println(err)
			}
			fmt.Println("")
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
