// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"encoding/binary"
	"fmt"
)

/*
Here is a diagram of the marshal format:

    0               1               2               3
    0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Marshal Version       |            Length...          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          ...Length            |             Flags             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |       p       |       p'      |        sparseLength...        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |       ...sparseLength         |bitsPerRegister|    Data...    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

const (
	marshalVersion    = 1
	marshalHeaderSize = 15

	marshalFlagSparse = 1
)

// Marshal serializes h into a byte slice that can be deserialized via
// Unmarshal. The data is naturally compressed, so don't bother trying
// to compress it any more.
func (h *HLLPP) Marshal() []byte {
	if h.sparse {
		h.flushTmpSet()
	}

	buf := make([]byte, marshalHeaderSize+len(h.data))

	offset := 0

	binary.BigEndian.PutUint16(buf[offset:], marshalVersion)
	offset += 2

	binary.BigEndian.PutUint32(buf[offset:], uint32(len(buf)))
	offset += 4

	var flags uint16
	if h.sparse {
		flags |= marshalFlagSparse
	}

	binary.BigEndian.PutUint16(buf[offset:], flags)
	offset += 2

	buf[offset] = h.p
	offset++

	buf[offset] = h.pp
	offset++

	binary.BigEndian.PutUint32(buf[offset:], h.sparseLength)
	offset += 4

	buf[offset] = byte(h.bitsPerRegister)
	offset++

	copy(buf[offset:], h.data)

	return buf
}

// Unmarshal deserializes a byte slice returned by Marshal back into an
// HLLPP object.
func Unmarshal(data []byte) (*HLLPP, error) {
	if len(data) < marshalHeaderSize {
		return nil, fmt.Errorf("data too short (%d bytes)", len(data))
	}

	offset := 0

	version := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	if version != marshalVersion {
		return nil, fmt.Errorf("unknown version: %d", version)
	}

	length := binary.BigEndian.Uint32(data[offset:])
	offset += 4

	if int(length) != len(data) {
		return nil, fmt.Errorf("length mismatch: header says %d, was %d", length, len(data))
	}

	flags := binary.BigEndian.Uint16(data[offset:])
	offset += 2

	p := data[offset]
	offset++

	pp := data[offset]
	offset++

	h, err := NewWithConfig(Config{
		Precision:       p,
		SparsePrecision: pp,
	})
	if err != nil {
		return nil, err
	}

	h.sparse = flags&marshalFlagSparse > 0

	h.sparseLength = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	h.bitsPerRegister = uint32(data[offset])
	offset++

	if len(data) > offset {
		h.data = make([]byte, len(data)-offset)
		copy(h.data, data[offset:])
	}

	return h, nil
}
