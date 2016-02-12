package storage

import "github.com/golang/protobuf/proto"

// CreateDom ...
const (
	CreateDom    = uint8(0)
	DeleteDom    = uint8(1)
	CreateSketch = uint8(2)
	DeleteSketch = uint8(3)
	Add          = uint8(4)
)

// Entry ...
type Entry struct {
	op  uint8
	msg proto.Message
	raw []byte
}

// Op ...
func (entry *Entry) Op() uint8 {
	return entry.op
}

// Args ...
func (entry *Entry) Msg() proto.Message {
	return entry.msg
}

// Args ...
func (entry *Entry) Raw() []byte {
	return entry.raw
}
