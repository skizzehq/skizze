package bloom

import (
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/sketches/wrappers/bloom/bloom"
)

//var logger = utils.GetLogger()

const defaultCapacity = 1000000

// Sketch is the toplevel Sketch to control the count-min-log implementation
type Sketch struct {
	*datamodel.Info
	impl *bloom.Filter
}

// NewSketch ...
func NewSketch(info *datamodel.Info) (*Sketch, error) {
	if info.Properties.Capacity == 0 {
		info.Properties.Capacity = defaultCapacity
	}
	sketch := bloom.New(uint(info.Properties.Capacity), 4)
	d := Sketch{info, sketch}
	return &d, nil
}

// Add ...
func (d *Sketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

// Marshal ...
func (d *Sketch) Marshal() ([]byte, error) {
	return d.impl.GobEncode()
}

// Get ...
func (d *Sketch) Get(data interface{}) (interface{}, error) {
	values := data.([][]byte)
	res := make(map[string]bool)
	for _, value := range values {
		res[string(value)] = d.impl.Test(value)
	}
	return res, nil
}

// Unmarshal ...
func (d *Sketch) Unmarshal(info *datamodel.Info, data []byte) error {
	sketch := &bloom.Filter{}
	err := sketch.GobDecode(data)

	if err != nil {
		return err
	}
	d.Info = info
	d.impl = sketch
	return nil
}
