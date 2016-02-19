package sketches

import (
	"github.com/seiflotfy/count-min-log"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// CMLSketch is the toplevel Sketch to control the count-min-log implementation
type CMLSketch struct {
	*datamodel.Info
	impl *cml.Sketch
}

// NewCMLSketch ...
func NewCMLSketch(info *datamodel.Info) (*CMLSketch, error) {
	sketch, err := cml.NewForCapacity16(uint64(info.Properties.GetMaxUniqueItems()), 0.01)
	d := CMLSketch{info, sketch}
	if err != nil {
		logger.Errorf("an error has occurred while saving CMLSketch: %s", err.Error())
	}
	return &d, nil
}

// Add ...
func (d *CMLSketch) Add(values [][]byte) (bool, error) {
	success := true
	for _, v := range values {
		if b := d.impl.IncreaseCount([]byte(v)); !b {
			success = false
		}
	}
	return success, nil
}

// Get ...
func (d *CMLSketch) Get(data interface{}) (interface{}, error) {
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
			Count: utils.Int64p(int64(d.impl.Frequency(v))),
		}
		tmpRes[string(v)] = res.Frequencies[i]
	}
	return res, nil
}
