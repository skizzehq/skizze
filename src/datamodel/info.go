package datamodel

import (
	pb "datamodel/protobuf"
	"fmt"
	"utils"
)

// Info represents a info string describing the sketch
type Info struct {
	*pb.Sketch
	locked bool
	id     string
}

// ID return a unique ID based on the name and type
func (info *Info) ID() string {
	if len(info.id) == 0 {
		info.id = fmt.Sprintf("%s.%s", info.GetName(), info.GetType())
	}
	return info.id
}

// Locked returns the lock state of the sketch
func (info *Info) Locked() bool {
	return info.locked
}

// Lock the Sketch
func (info *Info) Lock() {
	info.locked = true
}

// Unlock the Sketch
func (info *Info) Unlock() {
	info.locked = false
}

// Copy sketch
func (info *Info) Copy() *Info {
	typ := info.GetType()
	return &Info{
		Sketch: &pb.Sketch{
			Properties: &pb.SketchProperties{
				ErrorRate:      utils.Float32p(info.Properties.GetErrorRate()),
				MaxUniqueItems: utils.Int64p(info.Properties.GetMaxUniqueItems()),
				Size:           utils.Int64p(info.Properties.GetSize()),
			},
			State: &pb.SketchState{
				FillRate:     utils.Float32p(info.State.GetFillRate()),
				LastSnapshot: utils.Int64p(info.State.GetLastSnapshot()),
			},
			Name: utils.Stringp(info.GetName()),
			Type: &typ,
		},
	}
}

// NewEmptyProperties returns an empty property struct
func NewEmptyProperties() *pb.SketchProperties {
	return &pb.SketchProperties{
		ErrorRate:      utils.Float32p(0),
		MaxUniqueItems: utils.Int64p(0),
		Size:           utils.Int64p(0),
	}
}

// NewEmptyState returns an empty state struct
func NewEmptyState() *pb.SketchState {
	return &pb.SketchState{
		FillRate:     utils.Float32p(0),
		LastSnapshot: utils.Int64p(0),
	}
}

// NewEmptyInfo returns an empty info struct
func NewEmptyInfo() *Info {
	sketch := &pb.Sketch{
		Properties: NewEmptyProperties(),
		State:      NewEmptyState(),
	}
	return &Info{Sketch: sketch, locked: false}
}
