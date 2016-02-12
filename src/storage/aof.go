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
	in     chan *Entry
}

// NewAOF ...
func NewAOF(path string) *AOF {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	utils.PanicOnError(err)
	rdr := bufio.NewReader(file)
	wtr := bufio.NewWriter(file)
	in := make(chan *Entry)
	return &AOF{file, bufio.NewReadWriter(rdr, wtr), sync.RWMutex{}, in}
}

// Run ...
func (aof *AOF) Run() {
	go func() {
		for {
			select {
			case e := <-aof.in:
				aof.write(e)
			}
		}
	}()
}

func (aof *AOF) write(e *Entry) {
	line := fmt.Sprintf("%d|%s/", e.op, string(e.raw))
	if _, err := aof.buffer.WriteString(line); err != nil {
		fmt.Println(err)
	}
	if err := aof.buffer.Flush(); err != nil {
		fmt.Println(err)
	}
}

// Append ...
func (aof *AOF) Append(op uint8, msg proto.Message) error {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	e := &Entry{op, msg, raw}
	aof.in <- e
	return nil
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
	e := &Entry{uint8(op), nil, []byte(msg)}
	return e, nil
}
