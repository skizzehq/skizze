package topk

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/topk/go-topk"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

const defaultCapacity = 100.0

/*
Sketch is the toplevel sketch to control the HLL implementation
*/
type Sketch struct {
	*abstract.Info
	impl *topk.Stream
}

/*
ResultElement ...
*/
type ResultElement topk.Element

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	manager = storage.GetManager()
	err := manager.Create(info.ID)
	if err != nil {
		logger.Error.Println("an error has occurred while creating sketch: " + err.Error())
		return nil, err
	}
	if info.Properties["capacity"] == 0 {
		info.Properties["capacity"] = defaultCapacity
	}
	d := Sketch{info, topk.New(int(info.Properties["capacity"]))}

	return &d, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	str := string(value)
	err := d.impl.Insert(str, 1)
	if err != nil {
		logger.Error.Println("an error has occurred while populating sketch: " + err.Error())
		return false, err
	}
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Sketch) AddMultiple(values [][]byte) (bool, error) {

	for _, value := range values {
		str := string(value)
		err := d.impl.Insert(str, 1)
		if err != nil {
			logger.Error.Println("an error has occurred while populating sketch: " + err.Error())
			return false, err
		}
	}
	return true, nil
}

/*
Remove ...
*/
func (d *Sketch) Remove(value []byte) (bool, error) {
	logger.Error.Println("This sketch type does not support deletion")
	return false, errors.New("This sketch type does not support deletion")
}

/*
RemoveMultiple ...
*/
func (d *Sketch) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This sketch type does not support deletion")
	return false, errors.New("This sketch type does not support deletion")
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
	return true, nil
}

/*
GetType ...
*/
func (d *Sketch) GetType() string {
	return d.Type
}

/*
GetID ...
*/
func (d *Sketch) GetID() string {
	return d.ID
}

/*
GetFrequency ...
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	keys := d.impl.Keys()
	result := make([]ResultElement, len(keys), len(keys))
	for i, k := range keys {
		result[i] = ResultElement(k)
	}
	return result
}

/*
Marshal ...
*/
func (d *Sketch) Marshal() ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	// Encode (send) the value.
	err := enc.Encode(d.impl)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

/*
Unmarshal ...
*/
func Unmarshal(info *abstract.Info, data []byte) (*Sketch, error) {
	var network bytes.Buffer // Stand-in for a network connection
	_, err := network.Write(data)
	if err != nil {
		logger.Error.Println("an error has occurred while loading sketch from data: " + err.Error())
		return nil, err
	}
	dec := gob.NewDecoder(&network) // Will read from network.

	var counter topk.Stream
	err = dec.Decode(&counter)
	if err != nil {
		return nil, err
	}
	return &Sketch{info, &counter}, nil
}
