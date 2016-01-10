// Copyright (c) 2015, RetailNext, Inc.
// All rights reserved.

package hllpp

import (
	"crypto/sha1"
	"encoding/binary"
	"strings"
	"testing"
)

func TestMurmur(t *testing.T) {
	sixteen := "exactly-16-bytes"

	cases := []struct {
		input    []byte
		expected uint64
	}{
		{
			[]byte("\x00"),
			binary.LittleEndian.Uint64([]byte{181, 92, 255, 110, 229, 171, 16, 70, 131}),
		},
		{
			[]byte("\xFF"),
			binary.LittleEndian.Uint64([]byte{236, 144, 226, 164, 120, 55, 218, 71, 46}),
		},
		{
			[]byte("superlogical-skeleton"),
			binary.LittleEndian.Uint64([]byte{18, 139, 164, 68, 109, 45, 88, 40, 180}),
		},
		{
			[]byte(strings.Repeat("separator-quizzingly", 100)),
			binary.LittleEndian.Uint64([]byte{75, 88, 136, 147, 185, 104, 20, 147, 230}),
		},

		{
			[]byte("\x00"),
			binary.LittleEndian.Uint64([]byte{181, 92, 255, 110, 229, 171, 16, 70, 131}),
		},
		{
			[]byte("\x00\x01"),
			binary.LittleEndian.Uint64([]byte{76, 38, 171, 141, 197, 245, 179, 124, 68}),
		},
		{
			[]byte("\x00\x01\x02"),
			binary.LittleEndian.Uint64([]byte{190, 230, 83, 239, 47, 161, 114, 184, 182}),
		},
		{
			[]byte("\x00\x01\x02\x03"),
			binary.LittleEndian.Uint64([]byte{16, 175, 223, 13, 174, 148, 197, 225, 226}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04"),
			binary.LittleEndian.Uint64([]byte{54, 64, 249, 166, 212, 140, 238, 65, 246}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05"),
			binary.LittleEndian.Uint64([]byte{60, 4, 245, 164, 187, 58, 152, 102, 160}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06"),
			binary.LittleEndian.Uint64([]byte{104, 13, 75, 202, 135, 105, 76, 189, 135}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07"),
			binary.LittleEndian.Uint64([]byte{200, 47, 142, 214, 189, 225, 167, 71, 199}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08"),
			binary.LittleEndian.Uint64([]byte{50, 45, 129, 110, 15, 203, 180, 251, 185}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09"),
			binary.LittleEndian.Uint64([]byte{99, 228, 88, 158, 232, 37, 202, 207, 118}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A"),
			binary.LittleEndian.Uint64([]byte{136, 79, 86, 199, 71, 79, 123, 197, 176}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B"),
			binary.LittleEndian.Uint64([]byte{202, 165, 18, 146, 230, 167, 93, 179, 70}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C"),
			binary.LittleEndian.Uint64([]byte{194, 65, 95, 197, 242, 217, 82, 75, 252}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D"),
			binary.LittleEndian.Uint64([]byte{100, 109, 144, 53, 238, 51, 169, 95, 223}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E"),
			binary.LittleEndian.Uint64([]byte{233, 37, 73, 253, 152, 21, 35, 71, 233}),
		},
		{
			[]byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F"),
			binary.LittleEndian.Uint64([]byte{48, 63, 144, 145, 181, 36, 73, 68, 69}),
		},

		{
			[]byte(sixteen + "\x00"),
			binary.LittleEndian.Uint64([]byte{107, 80, 200, 85, 251, 132, 195, 106, 201}),
		},
		{
			[]byte(sixteen + "\x00\x01"),
			binary.LittleEndian.Uint64([]byte{8, 163, 120, 155, 17, 116, 0, 239, 159}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02"),
			binary.LittleEndian.Uint64([]byte{205, 145, 14, 223, 61, 213, 142, 251, 3}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03"),
			binary.LittleEndian.Uint64([]byte{165, 234, 67, 132, 62, 250, 205, 15, 107}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04"),
			binary.LittleEndian.Uint64([]byte{125, 248, 112, 82, 86, 0, 223, 50, 207}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05"),
			binary.LittleEndian.Uint64([]byte{228, 222, 242, 88, 206, 209, 208, 221, 22}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06"),
			binary.LittleEndian.Uint64([]byte{156, 172, 230, 68, 108, 255, 208, 237, 21}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07"),
			binary.LittleEndian.Uint64([]byte{0, 154, 196, 254, 211, 207, 108, 26, 61}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08"),
			binary.LittleEndian.Uint64([]byte{227, 158, 95, 195, 12, 14, 8, 67, 92}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09"),
			binary.LittleEndian.Uint64([]byte{63, 193, 84, 194, 216, 137, 100, 77, 231}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A"),
			binary.LittleEndian.Uint64([]byte{100, 22, 13, 55, 181, 18, 218, 178, 164}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B"),
			binary.LittleEndian.Uint64([]byte{164, 109, 195, 45, 7, 0, 15, 97, 174}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C"),
			binary.LittleEndian.Uint64([]byte{113, 136, 95, 250, 81, 190, 207, 152, 61}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D"),
			binary.LittleEndian.Uint64([]byte{191, 140, 212, 107, 211, 254, 58, 190, 108}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E"),
			binary.LittleEndian.Uint64([]byte{64, 129, 187, 169, 124, 88, 127, 169, 180}),
		},
		{
			[]byte(sixteen + "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F"),
			binary.LittleEndian.Uint64([]byte{162, 30, 109, 8, 38, 209, 12, 244, 183}),
		},
	}

	for i, c := range cases {
		if sum := murmurSum64(c.input); sum != c.expected {
			t.Errorf("#%d: got %d, expected %d", i, sum, c.expected)
		}
	}

	bigEndian = !bigEndian
	defer func() { bigEndian = !bigEndian }()
	for i, c := range cases {
		if sum := murmurSum64(c.input); sum != c.expected {
			t.Errorf("#%d: got %d, expected %d", i, sum, c.expected)
		}
	}
}

func BenchmarkMurmurSmall(b *testing.B) {
	data := []byte("zealotist")
	for i := 0; i < b.N; i++ {
		_ = murmurSum64(data)
	}
}

func BenchmarkMurmurMedium(b *testing.B) {
	data := []byte("zealotist-adventure crammer-empyromancy")
	for i := 0; i < b.N; i++ {
		_ = murmurSum64(data)
	}
}

func BenchmarkSHA1(b *testing.B) {
	h := sha1.New()
	data := []byte("zealotist")
	var tmp []byte
	for i := 0; i < b.N; i++ {
		_, err := h.Write(data)
		if err != nil {
			return
		}
		tmp = h.Sum(tmp[0:0])
		_ = binary.BigEndian.Uint64(tmp)
		h.Reset()
	}
}

func BenchmarkMurmurLarge(b *testing.B) {
	data := []byte(strings.Repeat("ergogram-artificer", 100))
	for i := 0; i < b.N; i++ {
		_ = murmurSum64(data)
	}
}
