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
	impl *hllpp.HLLPP
}

// NewHLLPPSketch ...
func NewHLLPPSketch(info *datamodel.Info) (*HLLPPSketch, error) {
	d := HLLPPSketch{info, hllpp.New()}
	return &d, nil
}

// Add ...
func (d *HLLPPSketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

// Get ...
func (d *HLLPPSketch) Get(interface{}) (interface{}, error) {
	return &pb.CardinalityResult{
		Cardinality: utils.Int64p(int64(d.impl.Count())),
	}, nil
}

// Marshal ...
func (d *HLLPPSketch) Marshal() ([]byte, error) {
	return d.impl.Marshal(), nil
}

// Unmarshal ...
func (d *HLLPPSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	impl, err := hllpp.Unmarshal(data)
	if err != nil {
		return err
	}
	d.impl = impl
	return nil
}
