package sketches

import (
	"github.com/skizzehq/count-min-log"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// CMLSketch is the toplevel Sketch to control the count-min-log implementation
type CMLSketch struct {
	*datamodel.Info
	impl      *cml.Sketch
	threshold *Dict
}

// NewCMLSketch ...
func NewCMLSketch(info *datamodel.Info) (*CMLSketch, error) {
	threshold := NewDict(info)
	d := CMLSketch{info, nil, threshold}
	return &d, nil
}

// Add ...
func (d *CMLSketch) Add(values [][]byte) (bool, error) {
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
			sketch, err := cml.NewForCapacity16(uint64(d.Info.Properties.GetMaxUniqueItems()), 0.01)
			if err != nil {
				logger.Errorf("an error has occurred while saving CMLSketch: %s", err.Error())
			}
			d.impl = sketch
		}
	}

	for _, v := range values {
		dict[string(v)]++
	}
	for v, count := range dict {
		if b := d.impl.BulkUpdate([]byte(v), count); !b {
			success = false
		}
	}
	return success, nil
}

// Get ...
func (d *CMLSketch) Get(data interface{}) (interface{}, error) {
	if d.threshold != nil {
		return d.threshold.Get(data)
	}

	values := data.([][]byte)
	res := &pb.FrequencyResult{
		Frequencies: make([]*pb.Frequency, len(values), len(values)),
	}
	tmpRes := make(map[string]*pb.Frequency)
	for i, v := range values {
		if r, ok := tmpRes[string(v)]; ok {
			res.Frequencies[i] = r
			continue
		}
		res.Frequencies[i] = &pb.Frequency{
			Value: utils.Stringp(string(v)),
			Count: utils.Int64p(int64(d.impl.Query(v))),
		}
		tmpRes[string(v)] = res.Frequencies[i]
	}
	return res, nil
}
