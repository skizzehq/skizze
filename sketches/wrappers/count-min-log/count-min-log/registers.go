package cml

import (
	"encoding/binary"

	"github.com/seiflotfy/skizze/storage"
)

type registers struct {
	d, w    uint
	storage *storage.ManagerStruct
	id      string
}

func newRegisters(id string, d, w uint) *registers {
	manager := storage.GetManager()
	return &registers{d, w, manager, id}
}

func (r *registers) get(i, j uint) (uint16, error) {
	newI := i * r.w * 2
	newJ := j * 2
	raw, err := r.storage.LoadData(r.id, int64(newI+newJ), 2)
	if err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint16(raw)
	return value, nil
}

func (r *registers) set(i, j uint, value uint16) error {
	newI := i * r.w * 2
	newJ := j * 2

	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, value)
	return r.storage.SaveData(r.id, data, int64(newI+newJ))
}

func (r *registers) save(data [][]uint16) error {
	bytes := make([]byte, r.d*r.w*2, r.d*r.w*2)
	for i := range data {
		for j, value := range data[i] {
			d := make([]byte, 2)
			pos := uint(i)*r.w*2 + uint(j)*2
			binary.LittleEndian.PutUint16(d, value)
			bytes[pos] = d[0]
			bytes[pos+1] = d[1]
		}
	}
	return r.storage.SaveData(r.id, bytes, 0)
}

func (r *registers) load() ([][]uint16, error) {
	raw, err := r.storage.LoadData(r.id, 0, int64(r.w*r.d*2))
	if err != nil {
		return nil, err
	}
	data := make([][]uint16, r.d, r.d)
	for i := range data {
		data[i] = make([]uint16, r.w, r.w)
		for j := range data[i] {
			pos := uint(i)*r.w*2 + uint(j)*2
			value := binary.LittleEndian.Uint16(raw[pos : pos+2])
			data[i][j] = value
		}
	}
	return data, nil
}

func (r *registers) reset() error {
	length := r.d * r.w * 2
	data := make([]byte, length, length)
	return r.storage.SaveData(r.id, data, 0)
}
