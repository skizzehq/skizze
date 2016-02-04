package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"utils"

	"github.com/golang/protobuf/proto"
)

// AOF ...
type AOF struct {
	file   *os.File
	buffer *bufio.ReadWriter
	lock   sync.RWMutex
}

// NewAOF ...
func NewAOF(path string) *AOF {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	utils.PanicOnError(err)
	rdr := bufio.NewReader(file)
	wtr := bufio.NewWriter(file)
	return &AOF{file, bufio.NewReadWriter(rdr, wtr), sync.RWMutex{}}
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

// AppendDomOp ...
func (aof *AOF) AppendDomOp(op uint8, dom proto.Message) error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	e, err := createEntry(dom)
	if err != nil {
		return err
	}
	if op != CreateDom && op != DeleteDom {
		return fmt.Errorf("No such op %d", op)
	}
	e.op = op
	return aof.write(e)
}

// AppendSketchOp ...
func (aof *AOF) AppendSketchOp(op uint8, sketch proto.Message) error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	e, err := createEntry(sketch)
	if err != nil {
		return err
	}
	e.op = op
	if op != CreateSketch && op != DeleteSketch {
		return fmt.Errorf("No such op %d", op)
	}
	return aof.write(e)
}

// AppendAddOp ...
func (aof *AOF) AppendAddOp(add proto.Message) error {
	aof.lock.Lock()
	defer aof.lock.Unlock()
	args, err := proto.Marshal(add)
	if err != nil {
		return err
	}
	e := &Entry{Add, args}
	return aof.write(e)
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
