// Copyright (c) 2015, RetailNext, Inc.
// This material contains trade secrets and confidential information of
// RetailNext, Inc.  Any use, reproduction, disclosure or dissemination
// is strictly prohibited without the explicit written permission
// of RetailNext, Inc.
// All rights reserved.
package hllpp

import "fmt"

func Example() {
	h := New()

	h.Add([]byte("barclay"))
	h.Add([]byte("reginald"))
	h.Add([]byte("barclay"))
	h.Add([]byte("broccoli"))

	fmt.Println(h.Count())
	// Output: 3
}

func ExampleNewWithConfig() {
	h, err := NewWithConfig(Config{
		Precision:       12,
		SparsePrecision: 14,
	})
	if err != nil {
		panic(err)
	}

	h.Add([]byte("qapla'"))
	h.Add([]byte("qapla'"))

	fmt.Println(h.Count())
	// Output: 1
}

func ExampleHLLPP_Marshal() {
	h := New()

	h.Add([]byte("hobbledehoyhood"))

	serialized := h.Marshal()
	_, err := Unmarshal(serialized)
	if err != nil {
		panic(err)
	}

	fmt.Println(h.Count())
	// Output: 1
}

func ExampleHLLPP_Merge() {
	h := New()
	h.Add([]byte("picard"))
	h.Add([]byte("janeway"))

	other := New()
	other.Add([]byte("picard"))
	other.Add([]byte("kirk"))

	err := h.Merge(other)
	if err != nil {
		panic(err)
	}

	fmt.Println(h.Count())
	// Output: 3
}
