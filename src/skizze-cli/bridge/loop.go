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

	"github.com/martinpinto/liner"
)

const helpString = `
skizze-cli %s

An interactive client to the Skizze database

Usage: skizze-cli [help]

Commands: (case insensitive)
  CREATE DOM  <name> <size> <maxUniqueItems>  Create a new Domain with options
  DESTROY DOM <name>                          Destroy a Domain

  CREATE CARD <name>                          Create a Cardinality Sketch
  CREATE MEMB <name>                          Create a Membership Sketch
  CREATE FREQ <name>                          Create a Frequency Sketch
  CREATE RANK <name>                          Create a Rankings Sketch

  LIST DOM                                    List existing Domains
  LIST                                        List existing Sketches

  INFO DOM <name>                             Get details of a Domain
  INFO <name>                                 Get details of a Sketch

  ADD DOM  <name> <value1> [value2...]        Add values to a Domain
  ADD FREQ <name> <value1> [value2...]        Add values to a frequency Sketch
  ADD MEMB <name> <value1> [value2...]        Add values to a membership Sketch
  ADD RANK <name> <value1> [value2...]        Add values to a rankings Sketch
  ADD CARD <name> <value1> [value2...]        Add values to a cardinality Sketch

  GET FREQ <name> <value1> [value2...]        Get the frequencies of the values in a FREQ Sketch
  GET MEMB <name> <value1> [value2...]        Get the memberships of the values in  a MEMB Sketch
  GET RANK <name>                             Get the top ranking values in a RANK Sketch
  GET CARD <name>                             Get the cardinality of a CARD Sketch

	QUIT                                        Exit skizze-cli

Shortcuts:
  Ctrl+d                                      Exit skizze-cli

Examples:
  CREATE DOM users 100 100000
  ADD DOM users neil seif martin conor neil conor seif seif seif
  GET FREQ users neil
  GET RANK users      
  GET CARD users

`

var (
	client     pb.SkizzeClient
	completion = []string{
		"create dom", "destroy dom",
		"create card", "create memb", "create freq", "create rank",
		"list", "list dom",
		"info", "info dom",
		"add dom", "add freq", "add memb", "add rank", "add card",
		"get freq", "get memb", "get rank", "get card",
		"help", "exit",
	}
	conn      *grpc.ClientConn
	historyFn = filepath.Join(os.TempDir(), ".skizze_history")
	w         = new(tabwriter.Writer)
	typeMap   = map[string]pb.SketchType{
		datamodel.HLLPP: pb.SketchType_CARD,
		datamodel.CML:   pb.SketchType_FREQ,
		datamodel.Bloom: pb.SketchType_MEMB,
		datamodel.TopK:  pb.SketchType_RANK,
	}
	version string
)

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

func evaluateQuery(query string) error {
	fields := getFields(query)
	if len(fields) != 0 && len(fields) <= 2 {
		//TODO: global stuff might be set
		switch strings.ToLower(fields[0]) {
		case "help":
			printHelp()
			return nil
		case "quit":
			tearDownClient(conn)
			os.Exit(0)
			return nil
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

func printHelp() {
	fmt.Printf(helpString, version)
}

// Run ...
func Run() {
	// Print help
	if len(os.Args) > 2 && (strings.HasSuffix(os.Args[1], "h") || strings.HasSuffix(os.Args[1], "help")) {
		printHelp()
		return
	}

	client, conn = setupClient()
	line := liner.NewLiner()
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	defer func() { _ = line.Close() }()

	line.SetCtrlCAborts(true)

	line.SetCompleter(func(line string) (c []string) {
		for _, n := range completion {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})

	if f, err := os.Open(historyFn); err == nil {
		if _, err := line.ReadHistory(f); err == nil {
			_ = f.Close()
		}
	}

	for {
		if query, err := line.Prompt("skizze> "); err == nil {
			if err := evaluateQuery(query); err != nil {
				log.Printf("Error evaluating query: %s", err.Error())
			}
			line.AppendHistory(query)
		} else if err == io.EOF {
			tearDownClient(conn)
			return
		} else if err == liner.ErrPromptAborted {
			fmt.Println("")
		} else {
			log.Printf("Error reading line: %s", err.Error())
		}

		if f, err := os.Create(historyFn); err != nil {
			log.Fatalf("Error writing history file: %s", err.Error())
		} else {
			if _, err := line.WriteHistory(f); err != nil {
				_ = f.Close()
			}
		}
	}
}
