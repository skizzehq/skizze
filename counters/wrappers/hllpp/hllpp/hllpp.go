// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

// Package hllpp implements the HyperLogLog++ cardinality estimator as specified
// in the HyperLogLog++ paper http://goo.gl/Z5Sqgu. hllpp uses a built-in
// non-streaming implementation of murmur3 to hash data as you add it to
// the estimator.
package hllpp

import (
	"errors"
	"fmt"
	"math"
)

// HLLPP represents a single HyperLogLog++ estimator. Create one via New().
// It is not safe to interact with an HLLPP object from multiple goroutines
// at once.
type HLLPP struct {
	// raw data be it sparse or dense (this makes serialization easier)
	data []byte

	// accumulates unsorted values in sparse mode
	tmpSet uint32Slice

	sparse       bool
	sparseLength uint32

	// how many bits we are using to store each register value
	bitsPerRegister uint32

	p uint8
	m uint32

	// p' and m'
	pp uint8
	mp uint32
}

// Approximate size in bytes of h (used for testing).
func (h *HLLPP) memSize() int {
	return cap(h.data) + 4*cap(h.tmpSet) + 20
}

// New creates a HyperLogLog++ estimator with p=14, p'=20.
func New() *HLLPP {
	h, err := NewWithConfig(Config{})
	if err != nil {
		panic(err)
	}
	return h
}

// Config is used to set configurable fields on a HyperLogLog++ via
// NewWithConfig.
type Config struct {
	// Precision (p). Must be in the range [4..16]. This value can be used
	// to adjust the typical relative error of the estimate. Space requirements
	// grow exponentially as this value is increased. Defaults to 14, the
	// recommended value, which gives an expected error of about 0.8%
	Precision uint8

	// Precision in sparse mode (p'). Must be in the range [p..25] for this
	// implementation. This value can be used to adjust the typical relative
	// error of the estimate when using the sparse representation (typically
	// for cardinalities below 8000 at p'=20). Lowering p' will allow the
	// estimator to remain in sparse mode longer, but will increase the relative
	// error. The HyperLogLog++ paper recommends 20 or 25. Defaults to 20 since
	// that still gives you a much lower error vs. p=14, but saves a signficant
	// amount of space vs. p'=25 (20-25% for cardinalities less than 5000).
	SparsePrecision uint8
}

// NewWithConfig creates a HyperLogLog++ estimator with the given Config.
func NewWithConfig(c Config) (*HLLPP, error) {
	if c.Precision == 0 {
		c.Precision = 14
	}

	if c.SparsePrecision == 0 {
		c.SparsePrecision = 20
	}

	p, pp := c.Precision, c.SparsePrecision
	if p < 4 || p > 16 || pp < p || pp > 25 {
		return nil, fmt.Errorf("invalid precision (p: %d, p': %d)", p, pp)
	}

	return &HLLPP{
		p:      p,
		pp:     pp,
		m:      1 << p,
		mp:     1 << pp,
		sparse: true,
	}, nil
}

// Add will hash v and add the result to the HyperLogLog++ estimator h. hllpp
// uses a built-in non-streaming implementation of murmur3.
func (h *HLLPP) Add(v []byte) {
	x := murmurSum64(v)

	if h.sparse {
		h.tmpSet = append(h.tmpSet, h.encodeHash(x))

		// is tmpSet >= 1/4 of memory limit?
		if 4*uint32(len(h.tmpSet))*8 >= 6*h.m/4 {
			h.flushTmpSet()
		}
	} else {
		idx := uint32(sliceBits64(x, 63, 64-h.p))
		rho := rho(x<<h.p | 1<<(h.p-1))
		h.updateRegisterIfBigger(idx, rho)
	}
}

func (h *HLLPP) updateRegisterIfBigger(idx uint32, rho uint8) {
	if rho > 31 && h.bitsPerRegister == 5 {
		h.bitsPerRegister = 6
		newData := make([]byte, h.m*h.bitsPerRegister/8)
		for i := uint32(0); i < h.m; i++ {
			setRegister(newData, 6, i, getRegister(h.data, 5, i))
		}
		h.data = newData
	}

	if rho > getRegister(h.data, h.bitsPerRegister, idx) {
		setRegister(h.data, h.bitsPerRegister, idx, rho)
	}
}

// Count returns the current cardinality estimate for h.
func (h *HLLPP) Count() uint64 {
	if h.sparse {
		h.flushTmpSet()
		return linearCounting(h.mp, h.mp-h.sparseLength)
	}

	var (
		est      float64
		numZeros uint32
	)
	for i := uint32(0); i < h.m; i++ {
		reg := getRegister(h.data, h.bitsPerRegister, i)
		est += 1.0 / float64(uint64(1)<<reg)
		if reg == 0 {
			numZeros++
		}
	}

	if numZeros > 0 {
		lc := linearCounting(h.m, numZeros)
		if lc < threshold[h.p-4] {
			return lc
		}
	}

	est = alpha(h.m) * float64(h.m) * float64(h.m) / est

	if est <= float64(h.m*5) {
		est -= h.estimateBias(est)
	}

	return uint64(est + 0.5)
}

// Merge turns h into the union of h and other. h and other must have the same
// p and p' values.
func (h *HLLPP) Merge(other *HLLPP) error {
	if h.p != other.p || h.pp != other.pp {
		return errors.New("HLLPPs have different parameters")
	}

	if h.sparse && !other.sparse {
		h.toNormal()
	}

	if other.sparse {
		other.flushTmpSet()
	}

	if h.sparse && other.sparse {
		tmpSet := make([]uint32, other.sparseLength)
		reader := newSparseReader(other.data)
		for index := 0; !reader.Done(); index++ {
			tmpSet[index] = reader.Next()
		}
		h.mergeSparse(tmpSet)
	} else if !h.sparse && !other.sparse {
		for i := uint32(0); i < h.m; i++ {
			rho := getRegister(other.data, other.bitsPerRegister, i)
			h.updateRegisterIfBigger(i, rho)
		}
	} else {
		reader := newSparseReader(other.data)
		for !reader.Done() {
			idx, rho := other.decodeHash(reader.Next(), other.p)
			h.updateRegisterIfBigger(idx, rho)
		}
	}

	return nil
}

func (h *HLLPP) toNormal() {
	if !h.sparse {
		return
	}

	if h.bitsPerRegister == 0 {
		h.bitsPerRegister = 5
	}

	newData := make([]byte, h.m*h.bitsPerRegister/8)

	reader := newSparseReader(h.data)
	for !reader.Done() {
		idx, rho := h.decodeHash(reader.Next(), h.p)

		if rho > 31 && h.bitsPerRegister == 5 {
			h.bitsPerRegister = 6
			h.toNormal()
			return
		}

		if rho > getRegister(newData, h.bitsPerRegister, idx) {
			setRegister(newData, h.bitsPerRegister, idx, rho)
		}
	}

	h.data = newData
	h.tmpSet = nil
	h.sparse = false
}

func linearCounting(m, v uint32) uint64 {
	return uint64(float64(m)*math.Log(float64(m)/float64(v)) + 0.5)
}

// slice out inclusive bit section [x.high..x.low]
func sliceBits64(x uint64, high, low uint8) uint64 {
	return (x << (63 - high)) >> (low + (63 - high))
}

// slice out inclusive bit section [x.high..x.low]
func sliceBits32(x uint32, high, low uint8) uint32 {
	return (x << (31 - high)) >> (low + (31 - high))
}

// number of leading zeros plus 1 (rho as in "ϱ" in paper)
func rho(x uint64) (z uint8) {
	for bit := uint64(1 << 63); bit&x == 0 && bit > 0; bit >>= 1 {
		z++
	}
	return z + 1
}
