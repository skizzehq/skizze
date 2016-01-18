package datamodel

// Sketcher ...
type Sketcher interface {
	Add([][]byte) (bool, error)
	Marshal() ([]byte, error)
	Get(interface{}) (interface{}, error)
	Unmarshal(*Info, []byte) error
}
