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

func (r *registers) reset() error {
	length := r.d * r.w * 2
	data := make([]byte, length, length)
	return r.storage.SaveData(r.id, data, 0)
}
