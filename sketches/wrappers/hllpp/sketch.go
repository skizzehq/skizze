package hllpp

import (
	"errors"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/hllpp/hllpp"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

/*
Sketch is the toplevel sketch to control the HLL implementation
*/
type Sketch struct {
	*abstract.Info
	impl *hllpp.HLLPP
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	d := Sketch{info, hllpp.New()}
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
	return uint(d.impl.Count())
}

/*
Clear ...
*/
func (d *Sketch) Clear() (bool, error) {
	return true, nil
}

/*
GetFrequency ...
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	return nil
}

/*
Marshal ...
*/
func (d *Sketch) Marshal() ([]byte, error) {
	return d.impl.Marshal(), nil
}

/*
Unmarshal ...
*/
func Unmarshal(info *abstract.Info, data []byte) (*Sketch, error) {
	counter, err := hllpp.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &Sketch{info, counter}, nil
}
