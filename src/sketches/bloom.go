package sketches

import (
	"github.com/willf/bloom"

	"datamodel"
)

// BloomSketch is the toplevel Sketch to control the count-min-log implementation
type BloomSketch struct {
	*datamodel.Info
	impl *bloom.BloomFilter
}

// NewBloomSketch ...
func NewBloomSketch(info *datamodel.Info) (*BloomSketch, error) {
	// FIXME: We are converting from int64 to uint
	sketch := bloom.New(uint(info.Properties.GetMaxUniqueItems()), 4)
	d := BloomSketch{info, sketch}
	return &d, nil
}

// Add ...
func (d *BloomSketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

// Marshal ...
func (d *BloomSketch) Marshal() ([]byte, error) {
	return d.impl.GobEncode()
}

// Get ...
func (d *BloomSketch) Get(data interface{}) (interface{}, error) {
	values := data.([][]byte)
	res := make([]*datamodel.Member, len(values), len(values))
	for i, v := range values {
		res[i] = &datamodel.Member{
			Key:    string(v),
			Member: d.impl.Test(v),
		}
	}
	return res, nil
}

// Unmarshal ...
func (d *BloomSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	sketch := &bloom.BloomFilter{}
	if err := sketch.GobDecode(data); err != nil {
		return err
	}
	d.Info = info
	d.impl = sketch
	return nil
}
