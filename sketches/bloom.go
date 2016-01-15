package sketches

import (
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/willf/bloom"
)

//var logger = utils.GetLogger()

const defaultCapacity = 1000000

// BloomSketch is the toplevel Sketch to control the count-min-log implementation
type BloomSketch struct {
	*datamodel.Info
	impl *bloom.BloomFilter
}

// NewBloomSketch ...
func NewBloomSketch(info *datamodel.Info) (*BloomSketch, error) {
	if info.Properties.Capacity == 0 {
		info.Properties.Capacity = defaultCapacity
	}
	sketch := bloom.New(uint(info.Properties.Capacity), 4)
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
