package cml

import "encoding/binary"

/*
Serialize ...
*/
func (sk *Sketch16) Serialize() ([]byte, error) {
	bytes := make([]byte, sk.k*sk.w*2, sk.k*sk.w*2)
	for i := range sk.store {
		for j, value := range sk.store[i] {
			d := make([]byte, 2)
			pos := uint(i)*sk.w*2 + uint(j)*2
			binary.LittleEndian.PutUint16(d, value)
			bytes[pos] = d[0]
			bytes[pos+1] = d[1]
		}
	}
	return bytes, nil
}

/*
Deserialize ...
*/
func (sk *Sketch16) Deserialize(raw []byte) ([][]uint16, error) {
	data := make([][]uint16, sk.k, sk.k)
	for i := range data {
		data[i] = make([]uint16, sk.w, sk.w)
		for j := range data[i] {
			pos := uint(i)*sk.w*2 + uint(j)*2
			value := binary.LittleEndian.Uint16(raw[pos : pos+2])
			data[i][j] = value
		}
	}
	return data, nil
}

/*
SetRegisters ...
*/
func (sk *Sketch16) SetRegisters(store [][]uint16) {
	sk.store = store
}
