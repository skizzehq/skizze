// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"testing"
)

func TestSparseReaderWriter(t *testing.T) {
	writer := newSparseWriter()

	if writer.Len() != 0 {
		t.Errorf("got %d", writer.Len())
	}

	if len(writer.Bytes()) != 0 {
		t.Errorf("got %+v", writer.Bytes())
	}

	reader := newSparseReader(writer.Bytes())

	if !reader.Done() {
		t.Errorf("should be done")
	}

	writer.Append(127, 0, 1)
	// same idx, but bigger rho than previous, use this one
	writer.Append(126, 0, 2)

	if writer.Len() != 1 || len(writer.Bytes()) != 1 {
		t.Errorf("got %d", writer.Len())
	}

	// show we are storing deltas since 128 takes two bytes as
	// a varint
	writer.Append(128, 1, 0)
	if writer.Len() != 2 || len(writer.Bytes()) != 2 {
		t.Errorf("got %d", writer.Len())
	}

	reader = newSparseReader(writer.Bytes())

	if reader.Done() {
		t.Errorf("shouldn't be done")
	}

	if reader.Peek() != 126 {
		t.Errorf("got %d", reader.Peek())
	}

	if reader.Peek() != 126 {
		t.Errorf("got %d", reader.Peek())
	}

	reader.Advance()

	if reader.Done() {
		t.Errorf("shouldn't be done")
	}

	if reader.Peek() != 128 {
		t.Errorf("got %d", reader.Peek())
	}

	reader.Advance()

	if !reader.Done() {
		t.Errorf("should be done")
	}
}
