package datamodel

// Sketcher ...
type Sketcher interface {
	Add([][]byte) (bool, error)
	Get(interface{}) (interface{}, error)
}
