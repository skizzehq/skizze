package datamodel

import "fmt"

// Properties ...
type Properties struct {
	Capacity uint `json:"capacity"`
	Rank     uint `json:"rank"`
}

// State represents a info string describing the sketch
type State struct {
	locked                bool
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

// Locked returns the lock state of the sketch
func (info *Info) Locked() bool {
	return info.State.locked
}

// Lock the Sketch
func (info *Info) Lock() {
	info.State.locked = true
}

// Unlock the Sketch
func (info *Info) Unlock() {
	info.State.locked = false
}

// NewEmptyProperties returns an empty property struct
func NewEmptyProperties() *Properties {
	return &Properties{}
}

// NewEmptyState returns an empty state struct
func NewEmptyState() *State {
	return &State{
		locked:                false,
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
