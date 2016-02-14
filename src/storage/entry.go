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

// OpType ...
func (entry *Entry) OpType() uint8 {
	return entry.op
}

// Msg ...
func (entry *Entry) Msg() proto.Message {
	return entry.msg
}

// RawMsg ...
func (entry *Entry) RawMsg() []byte {
	return entry.raw
}
