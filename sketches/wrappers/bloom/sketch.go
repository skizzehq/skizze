package bloom

import (
	"errors"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/bloom/bloom"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

/*
Sketch is the toplevel Sketch to control the count-min-log implementation
*/
type Sketch struct {
	*abstract.Info
	impl *bloom.Filter
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	// standard size
	sketch := bloom.New(1000000, 4)
	d := Sketch{info, sketch}
	return &d, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.impl.Add(value)
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Sketch) AddMultiple(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

/*
Remove ...
*/
func (d *Sketch) Remove(value []byte) (bool, error) {
	logger.Error.Println("This Sketch type does not support deletion")
	return false, errors.New("This Sketch type does not support deletion")
}

/*
RemoveMultiple ...
*/
func (d *Sketch) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This Sketch type does not support deletion")
	return false, errors.New("This Sketch type does not support deletion")
}

/*
GetCount ...
*/
func (d *Sketch) GetCount() uint {
	return 0
}

/*
Clear ...
*/
func (d *Sketch) Clear() (bool, error) {
	d.impl.ClearAll()
	return true, nil
}

/*
Marshal ...
*/
func (d *Sketch) Marshal() ([]byte, error) {
	return d.impl.GobEncode()
}

/*
GetFrequency ...
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	res := make(map[string]bool)
	for _, value := range values {
		res[string(value)] = d.impl.Test(value)
	}
	return res
}

/*
Unmarshal ...
*/
func Unmarshal(info *abstract.Info, data []byte) (*Sketch, error) {
	sketch := &bloom.Filter{}
	err := sketch.GobDecode(data)

	if err != nil {
		return nil, err
	}
	return &Sketch{info, sketch}, nil
}
