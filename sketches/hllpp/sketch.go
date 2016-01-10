package hllpp

import (
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/sketches/hllpp/hllpp"
)

//var logger = utils.GetLogger()

// Sketch is the toplevel sketch to control the HLL implementation
type Sketch struct {
	*datamodel.Info
	impl *hllpp.HLLPP
}

// NewSketch ...
func NewSketch(info *datamodel.Info) (*Sketch, error) {
	d := Sketch{info, hllpp.New()}
	return &d, nil
}

// Add ...
func (d *Sketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

// Get ...
func (d *Sketch) Get(interface{}) (interface{}, error) {
	return uint(d.impl.Count()), nil
}

// Marshal ...
func (d *Sketch) Marshal() ([]byte, error) {
	return d.impl.Marshal(), nil
}

// Unmarshal ...
func (d *Sketch) Unmarshal(info *datamodel.Info, data []byte) error {
	impl, err := hllpp.Unmarshal(data)
	if err != nil {
		return err
	}
	d.impl = impl
	return nil
}
