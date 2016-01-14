package sketches

import (
	"github.com/seiflotfy/count-min-log"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

//CMLSketch is the toplevel Sketch to control the count-min-log implementation
type CMLSketch struct {
	*datamodel.Info
	impl *cml.Sketch
}

//NewCMLSketch ...
func NewCMLSketch(info *datamodel.Info) (*CMLSketch, error) {
	if info.Properties.Capacity == 0 {
		info.Properties.Capacity = defaultCapacity
	}
	sketch, err := cml.NewForCapacity16(uint64(info.Properties.Capacity), 0.01)
	d := CMLSketch{info, sketch}
	if err != nil {
		logger.Error.Printf("an error has occurred while saving CMLSketch: %s", err.Error())
	}
	return &d, nil
}

//Add ...
func (d *CMLSketch) Add(values [][]byte) (bool, error) {
	success := true
	for _, v := range values {
		if b := d.impl.IncreaseCount([]byte(v)); !b {
			success = false
		}
	}
	return success, nil
}

//Get ...
func (d *CMLSketch) Get(data interface{}) (interface{}, error) {
	values := data.([][]byte)
	res := make(map[string]uint)
	for _, value := range values {
		count := d.impl.Frequency(value)
		res[string(value)] = uint(count)
	}
	return res, nil
}

//Marshal ...
func (d *CMLSketch) Marshal() ([]byte, error) {
	return d.impl.Marshal()
}

// Unmarshal ...
func (d *CMLSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	impl, err := cml.Unmarshal(data)
	if err != nil {
		return err
	}
	d.impl = impl
	return nil
}
