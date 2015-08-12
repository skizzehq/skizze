// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

// Non-streaming implementation of murmur3. Not only is it faster to begin with
// vs cryptographic hashing functions, but I can avoid memory allocations by
// not using the go streaming hash.Hash interface.

const (
	murmurC1 = 0x87c37b91114253d5
	murmurC2 = 0x4cf5ad432745937f
)

var bigEndian bool

func init() {
	var t uint16 = 1
	bigEndian = (*[2]byte)(unsafe.Pointer(&t))[0] == 0
}

// This is a port of MurmurHash3_x64_128 from MurmurHash3.cpp
func murmurSum64(data []byte) uint64 {
	var h1, h2, k1, k2 uint64

	len := len(data)

	nBlocks := len / 16

	var data64 []uint64

	if bigEndian {
		data64 = make([]uint64, nBlocks*2)
		for i := 0; i < nBlocks*2; i++ {
			data64[i] = binary.LittleEndian.Uint64(data[8*i:])
		}
	} else {
		dataHeader := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		data64Header := (*reflect.SliceHeader)(unsafe.Pointer(&data64))
		data64Header.Data = dataHeader.Data
		data64Header.Len = 2 * nBlocks
		data64Header.Cap = 2 * nBlocks
	}

	for i := 0; i < nBlocks; i++ {
		k1 = data64[2*i]
		k2 = data64[2*i+1]

		k1 *= murmurC1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= murmurC2
		h1 ^= k1

		h1 = (h1 << 27) | (h1 >> (64 - 27))
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= murmurC2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= murmurC1
		h2 ^= k2

		h2 = (h2 << 31) | (h2 >> (64 - 31))
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	tail := data[nBlocks*16:]

	k1, k2 = 0, 0

	switch len & 15 {
	case 15:
		k2 ^= uint64(tail[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(tail[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(tail[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(tail[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(tail[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(tail[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(tail[8]) << 0
		k2 *= murmurC2
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= murmurC1
		h2 ^= k2
		fallthrough
	case 8:
		k1 ^= uint64(tail[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(tail[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(tail[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(tail[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(tail[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(tail[0]) << 0
		k1 *= murmurC1
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= murmurC2
		h1 ^= k1
	}

	h1 ^= uint64(len)
	h2 ^= uint64(len)

	h1 += h2
	h2 += h1

	h1 ^= h1 >> 33
	h1 *= 0xff51afd7ed558ccd
	h1 ^= h1 >> 33
	h1 *= 0xc4ceb9fe1a85ec53
	h1 ^= h1 >> 33

	h2 ^= h2 >> 33
	h2 *= 0xff51afd7ed558ccd
	h2 ^= h2 >> 33
	h2 *= 0xc4ceb9fe1a85ec53
	h2 ^= h2 >> 33

	h1 += h2
	h2 += h1

	return h1
}
