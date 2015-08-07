// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import "sort"

// create a mask of numOnes 1's, shifted left shift bits
func mask(numOnes, shift uint32) uint32 {
	return ((1 << numOnes) - 1) << shift
}

func setRegister(data []byte, bitsPerRegister, idx uint32, rho uint8) {
	bitIdx := idx * bitsPerRegister
	byteOffset := bitIdx / 8
	bitOffset := bitIdx % 8

	if 8-bitOffset >= bitsPerRegister {
		// can all fit in first byte

		leftShift := 8 - bitsPerRegister - bitOffset

		// clear existing register value
		data[byteOffset] &= ^byte(mask(bitsPerRegister, leftShift))
		data[byteOffset] |= rho << leftShift
	} else {
		// spread over two bytes

		numBitsInFirstByte := bitsPerRegister - (8 - bitOffset)

		data[byteOffset] &= ^byte(mask(8-bitOffset, 0))
		data[byteOffset] |= rho >> numBitsInFirstByte

		data[byteOffset+1] &= ^byte(mask(numBitsInFirstByte, 8-numBitsInFirstByte))
		data[byteOffset+1] |= rho << (8 - numBitsInFirstByte)
	}
}

func getRegister(data []byte, bitsPerRegister, idx uint32) uint8 {
	bitIdx := idx * bitsPerRegister
	byteOffset := bitIdx / 8
	bitOffset := bitIdx % 8

	if 8-bitOffset >= bitsPerRegister {
		// all fit in first byte
		return (data[byteOffset] >> (8 - bitOffset - bitsPerRegister)) & byte(mask(bitsPerRegister, 0))
	}
	// spread over two bytes

	numBitsInFirstByte := bitsPerRegister - (8 - bitOffset)

	rho := data[byteOffset] << numBitsInFirstByte
	rho |= data[byteOffset+1] >> (8 - numBitsInFirstByte)
	return rho & byte(mask(bitsPerRegister, 0))

}

func alpha(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}

func (h *HLLPP) estimateBias(e float64) float64 {
	estimates := rawEstimateData[h.p-4]
	biases := biasData[h.p-4]

	index := sort.SearchFloat64s(estimates, e)

	if index == 0 {
		return biases[0]
	} else if index == len(estimates) {
		return biases[len(biases)-1]
	}

	e1, e2 := estimates[index-1], estimates[index]
	b1, b2 := biases[index-1], biases[index]

	r := (e - e1) / (e2 - e1)
	return b1*(1-r) + b2*r
}
