// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"encoding/binary"
	"sort"
)

type uint32Slice []uint32

func (s uint32Slice) Len() int {
	return len(s)
}

func (s uint32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s uint32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sparseReader struct {
	data    []byte
	idx     int
	lastVal uint32
	lastN   int
}

func newSparseReader(data []byte) *sparseReader {
	return &sparseReader{data: data}
}

func (iter *sparseReader) Advance() {
	iter.idx += iter.lastN
	iter.lastN = 0
}

func (iter *sparseReader) Peek() uint32 {
	if iter.lastN > 0 {
		return iter.lastVal
	}

	v, n := binary.Uvarint(iter.data[iter.idx:])

	v32 := uint32(v)
	v32 += iter.lastVal

	iter.lastN = n
	iter.lastVal = v32

	return v32
}

func (iter *sparseReader) Next() uint32 {
	v := iter.Peek()
	iter.Advance()
	return v
}

func (iter *sparseReader) Done() bool {
	return iter.idx >= len(iter.data)
}

type sparseWriter struct {
	data []byte

	lastVal uint32

	hasCurrVal bool
	currVal    uint32
	currIdx    uint32
	currRho    uint8

	varIntBuf []byte
	length    uint32
}

// This takes the index and rho well so it can easily discard duplicate indexes
// and pick the biggest rho among the duplicates since tmpSet is sorted by index
// but not by rho.
func (writer *sparseWriter) Append(k, idx uint32, rho uint8) {
	if writer.hasCurrVal {
		if idx == writer.currIdx {
			if rho > writer.currRho {
				writer.currRho = rho
				writer.currVal = k
			}
			return
		}
		writer.commit()
	}

	writer.hasCurrVal = true
	writer.currVal = k
	writer.currIdx = idx
	writer.currRho = rho
}

func (writer *sparseWriter) commit() {
	n := binary.PutUvarint(writer.varIntBuf, uint64(writer.currVal-writer.lastVal))
	writer.data = append(writer.data, writer.varIntBuf[:n]...)
	writer.lastVal = writer.currVal
	writer.length++
	writer.hasCurrVal = false
}

func (writer *sparseWriter) Bytes() []byte {
	if writer.hasCurrVal {
		writer.commit()
	}
	return writer.data
}

func (writer *sparseWriter) Len() uint32 {
	if writer.hasCurrVal {
		writer.commit()
	}
	return writer.length
}

func newSparseWriter() *sparseWriter {
	return &sparseWriter{
		varIntBuf: make([]byte, binary.MaxVarintLen32),
	}
}

func (h *HLLPP) flushTmpSet() {
	if len(h.tmpSet) == 0 {
		return
	}

	sort.Sort(h.tmpSet)
	h.mergeSparse(h.tmpSet)
	h.tmpSet = nil
}

func (h *HLLPP) mergeSparse(tmpSet []uint32) {

	iter := newSparseReader(h.data)
	writer := newSparseWriter()

	var tmpI int

	// deduping by index and choosing biggest rho is handled in the writer

	for !iter.Done() || tmpI < len(tmpSet) {
		if iter.Done() {
			idx, rho := h.decodeHash(tmpSet[tmpI], h.pp)
			writer.Append(tmpSet[tmpI], idx, rho)
			tmpI++
			continue
		}

		sparseVal := iter.Peek()
		sparseIdx, sparseR := h.decodeHash(sparseVal, h.pp)

		if tmpI == len(tmpSet) {
			writer.Append(sparseVal, sparseIdx, sparseR)
			iter.Advance()
			continue
		}

		tmpVal := tmpSet[tmpI]
		tmpIdx, tmpR := h.decodeHash(tmpVal, h.pp)

		if sparseIdx < tmpIdx {
			writer.Append(sparseVal, sparseIdx, sparseR)
			iter.Advance()
		} else if sparseIdx > tmpIdx {
			writer.Append(tmpVal, tmpIdx, tmpR)
			tmpI++
		} else {
			if sparseR > tmpR {
				writer.Append(sparseVal, sparseIdx, sparseR)
			} else {
				writer.Append(tmpVal, tmpIdx, tmpR)
			}
			iter.Advance()
			tmpI++
		}
	}

	h.data = writer.Bytes()
	h.sparseLength = writer.Len()

	// is sparse data bigger than dense data would be?
	if uint32(len(h.data))*8 >= 6*h.m {
		h.toNormal()
	}
}

func (h *HLLPP) encodeHash(x uint64) uint32 {
	if sliceBits64(x, 63-h.p, 64-h.pp) == 0 {
		r := rho((sliceBits64(x, 63-h.pp, 0) << h.pp) | (1<<h.pp - 1))
		return uint32(sliceBits64(x, 63, 64-h.pp)<<7 | uint64(r<<1) | 1)
	}

	return uint32(sliceBits64(x, 63, 64-h.pp) << 1)
}

// Return index with respect to "p" arg, and rho with respect to h.p. This is so
// the h.pp index can be recovered easily when flushing the tmpSet.
func (h *HLLPP) decodeHash(k uint32, p uint8) (_ uint32, r uint8) {
	if k&1 > 0 {
		r = uint8(sliceBits32(k, 6, 1)) + (h.pp - h.p)
	} else {
		r = rho((uint64(k) | 1) << (64 - (h.pp + 1) + h.p))
	}

	return h.getIndex(k, p), r
}

// Return index with respect to precision "p".
func (h *HLLPP) getIndex(k uint32, p uint8) uint32 {
	if k&1 > 0 {
		return sliceBits32(k, 6+h.pp, 1+6+h.pp-p)
	}
	return sliceBits32(k, h.pp, 1+h.pp-p)
}
