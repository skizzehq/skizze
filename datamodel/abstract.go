package datamodel

import "fmt"

/*
HLLPP	=> HyperLogLogPlusPlus
CML		=> Count-min-log sketch
TopK	=> Top-K
Dict  => dictionary
Bloom => Bloom Filter
*/
const (
	HLLPP = "card"
	CML   = "freq"
	TopK  = "rank"
	Dict  = "dict"
	Bloom = "memb"
)

// Element ...
type Element struct {
	Key   string
	Count int
	Error int
}

// Sketch ...
type Sketch interface {
	Add([][]byte) (bool, error)
	Marshal() ([]byte, error)
	Get(interface{}) (interface{}, error)
	Unmarshal(*Info, []byte) error
}

// Properties ...
type Properties struct {
	Capacity uint `json:"capacity"`
}

// State represents a info string describing the sketch
type State struct {
	Additions             uint `json:"adds"`
	LastSnapshotTimestamp uint `json:"last_snapshot_timestamp"`
}

// Info represents a info string describing the sketch
type Info struct {
	id         string
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	State      *State      `json:"state"`
	Properties *Properties `json:"properties"`
}

// ID return a unique ID based on the name and type
func (info *Info) ID() string {
	if len(info.id) == 0 {
		info.id = fmt.Sprintf("%s.%s", info.Name, info.Type)
	}
	return info.id
}

// NewEmptyProperties returns an empty property struct
func NewEmptyProperties() *Properties {
	return &Properties{}
}

// NewEmptyState returns an empty state struct
func NewEmptyState() *State {
	return &State{
		Additions:             0,
		LastSnapshotTimestamp: 0,
	}
}

// NewEmptyInfo returns an empty info struct
func NewEmptyInfo() *Info {
	return &Info{
		Properties: NewEmptyProperties(),
		State:      NewEmptyState(),
	}
}
