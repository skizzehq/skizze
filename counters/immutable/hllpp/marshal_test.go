// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"fmt"
	"reflect"
	"testing"
)

func hllpEqual(h1, h2 HLLPP) bool {
	return reflect.DeepEqual(h1, h2)
}

func marshalUnmarshal(h *HLLPP) error {
	unmarshaled, err := Unmarshal(h.Marshal())
	if err != nil {
		panic(err)
	}

	if unmarshaled.Count() != h.Count() {
		return fmt.Errorf("mismatched count: got %d, expected %d", unmarshaled.Count(), h.Count())
	}

	if !hllpEqual(*h, *unmarshaled) {
		return fmt.Errorf("got %+v, expected %+v", unmarshaled, h)
	}
	return nil
}

func TestMarshal(t *testing.T) {
	h := New()

	if err := marshalUnmarshal(h); err != nil {
		t.Error(err)
	}

	h.Add(intToBytes(1))

	if err := marshalUnmarshal(h); err != nil {
		t.Error(err)
	}

	for i := uint64(0); i < 1000; i++ {
		h.Add(intToBytes(i))
	}

	if !h.sparse {
		t.Error("Expecting sparse")
	}

	if err := marshalUnmarshal(h); err != nil {
		t.Error(err)
	}

	for i := uint64(0); i < 100000; i++ {
		h.Add(intToBytes(i))
	}

	if h.sparse {
		t.Error("Expecting dense")
	}

	if err := marshalUnmarshal(h); err != nil {
		t.Error(err)
	}
}

func TestUnmarshalErrors(t *testing.T) {
	uh, err := Unmarshal(nil)
	if uh != nil || err == nil {
		t.Error("Expected nil hll and some error")
	}

	uh, err = Unmarshal([]byte{})
	if uh != nil || err == nil {
		t.Error("Expected nil hll and some error")
	}

	h := New()
	for i := uint64(0); i < 10000; i++ {
		h.Add(intToBytes(i))
	}
	uh, err = Unmarshal(h.Marshal()[0:100])
	if uh != nil || err == nil {
		t.Error("Expected nil hll and some error")
	}
}
