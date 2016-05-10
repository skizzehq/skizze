package sketches

import (
	"github.com/retailnext/hllpp"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// HLLPPSketch is the toplevel sketch to control the HLL implementation
type HLLPPSketch struct {
	*datamodel.Info
	impl      *hllpp.HLLPP
	threshold *Dict
}

// NewHLLPPSketch ...
func NewHLLPPSketch(info *datamodel.Info) (*HLLPPSketch, error) {
	threshold := NewDict(info)
	d := HLLPPSketch{info, nil, threshold}
	return &d, nil
}

// Add ...
func (d *HLLPPSketch) Add(values [][]byte) (bool, error) {
	success := true
	dict := make(map[string]uint)
	if d.threshold != nil {
		s, err := d.threshold.Add(values)
		success = s
		if err != nil {
			return false, err
		}
		if !d.threshold.IsFull() {
			return true, nil
		}
		values = d.threshold.Keys()
		d.threshold = nil
		if d.impl == nil {
			sketch := hllpp.New()
			d.impl = sketch
		}
	}

	for _, v := range values {
		dict[string(v)]++
	}
	for v := range dict {
		d.impl.Add([]byte(v))
	}
	return success, nil
}

// Get ...
func (d *HLLPPSketch) Get(interface{}) (interface{}, error) {
	if d.threshold != nil {
		return d.threshold.Get(nil)
	}
	return &pb.CardinalityResult{
		Cardinality: utils.Int64p(int64(d.impl.Count())),
	}, nil
}
