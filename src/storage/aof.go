package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"utils"

	"github.com/golang/protobuf/proto"
)

// AOF ...
type AOF struct {
	file   *os.File
	buffer *bufio.ReadWriter

	// writeChan takes entries, which are written to file
	writeChan chan *Entry
	// writeState gets set to false if writeChan is close to prevent panics
	// when trying to put a value on a closed chan
	writeState bool
	// writtenRequest holds all successful requests
	writtenChan chan *Entry
}

// NewAOF ...
func NewAOF(path string) *AOF {
	// Open or create file and create ReadWriter
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	utils.PanicOnError(err)
	rwtr := bufio.NewReadWriter(
		bufio.NewReader(file),
		bufio.NewWriter(file),
	)

	// Create AOF
	aof := &AOF{
		file:        file,
		buffer:      rwtr,
		writeChan:   make(chan *Entry),
		writeState:  true,
		writtenChan: make(chan *Entry),
	}

	return aof
}

func (aof *AOF) write(e *Entry) error {
	line := fmt.Sprintf("%d|%s/", e.op, string(e.args))
	_, err := aof.buffer.WriteString(line)
	if err != nil {
		return err
	}
	err = aof.buffer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func createEntry(dom proto.Message) (*Entry, error) {
	args, err := proto.Marshal(dom)
	if err != nil {
		return nil, err
	}
	return &Entry{args: args}, nil
}

func (aof *AOF) queueWrite(e *Entry) {
	if aof.writeState {
		aof.writeChan <- e
	}
}

// Run reads from the writeChan and executes write requests in order
// If the writeChan is closed, this method exits as well
func (aof *AOF) Run() {
	for {
		if e, more := <-aof.writeChan; more {
			if err := aof.write(e); err == nil {
				select {
				case aof.writtenChan <- e:
				default:
				}
			}
		} else {
			// Channel closed, bye
			break
		}
	}
}

// WrittenChan returns the channel for written entries ro
func (aof *AOF) WrittenChan() <-chan *Entry {
	return aof.writtenChan
}

// AppendDomOp ...
func (aof *AOF) AppendDomOp(op uint8, dom proto.Message) {
	e, err := createEntry(dom)
	if err != nil {
		return // err
	}
	if op != CreateDom && op != DeleteDom {
		return // fmt.Errorf("No such op %d", op)
	}
	e.op = op

	aof.queueWrite(e)
}

// AppendSketchOp ...
func (aof *AOF) AppendSketchOp(op uint8, sketch proto.Message) {
	e, err := createEntry(sketch)
	if err != nil {
		return // err
	}
	e.op = op
	if op != CreateSketch && op != DeleteSketch {
		return // fmt.Errorf("No such op %d", op)
	}
	aof.queueWrite(e)
}

// AppendAddOp ...
func (aof *AOF) AppendAddOp(add proto.Message) {
	args, err := proto.Marshal(add)
	if err != nil {
		return // err
	}
	e := &Entry{Add, args}
	aof.queueWrite(e)
}

// Read ...
func (aof *AOF) Read() (*Entry, error) {
	line, err := aof.buffer.ReadBytes('/')
	if err != nil {
		return nil, err
	}
	line = line[:len(line)-1]
	rs := strings.Split(string(line), "|")
	msg := rs[1]
	op, err := strconv.Atoi(rs[0])
	if err != nil {
		return nil, err
	}
	e := &Entry{uint8(op), []byte(msg)}
	return e, nil
}

// Stop stops writing and waits until all write actions are finished
func (aof *AOF) Stop() {
	aof.writeState = false
	close(aof.writeChan)
}
