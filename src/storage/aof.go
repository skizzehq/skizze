package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"utils"

	"github.com/golang/protobuf/proto"
	"github.com/njpatel/loggo"
)

var logger = loggo.GetLogger("storage")

// AOF ...
type AOF struct {
	file     *os.File
	buffer   *bufio.ReadWriter
	lock     sync.RWMutex
	inChan   chan *Entry
	tickChan <-chan time.Time
}

// NewAOF ...
func NewAOF(path string) *AOF {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	utils.PanicOnError(err)
	rdr := bufio.NewReader(file)
	wtr := bufio.NewWriter(file)
	inChan := make(chan *Entry, 100)
	tickChan := time.NewTicker(time.Second).C
	return &AOF{
		file:     file,
		buffer:   bufio.NewReadWriter(rdr, wtr),
		lock:     sync.RWMutex{},
		inChan:   inChan,
		tickChan: tickChan,
	}
}

// Run ...
func (aof *AOF) Run() {
	go func() {
		for {
			select {
			case e := <-aof.inChan:
				aof.write(e)
			case <-aof.tickChan:
				if err := aof.buffer.Flush(); err != nil {
					logger.Errorf("an error has occurred while flushing AOF: %s", err.Error())
				}
			}
		}
	}()
}

func (aof *AOF) write(e *Entry) {
	line := fmt.Sprintf("%d|%s/", e.op, string(e.raw))
	if _, err := aof.buffer.WriteString(line); err != nil {
		logger.Errorf("an error has ocurred while writing AOF: %s", err.Error())
	}
}

// Append ...
func (aof *AOF) Append(op uint8, msg proto.Message) error {
	raw, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	e := &Entry{op, msg, raw}
	aof.inChan <- e
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
