package cml

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/count-min-log/count-min-log"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

const defaultEpsilon = 0.000003397855
const defaultDelta = 0.99

/*
Sketch is the toplevel Sketch to control the count-min-log implementation
*/
type Sketch struct {
	*abstract.Info
	impl *cml.Sketch16
	lock sync.RWMutex
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	epsilon := 0.0
	if eps, ok := info.Properties["epsilon"]; ok {
		epsilon = eps
	} else {
		epsilon = defaultEpsilon
		info.Properties["epsilon"] = epsilon
	}

	delta := 0.0
	if d, ok := info.Properties["delta"]; ok {
		delta = d
	} else {
		delta = defaultDelta
		info.Properties["delta"] = delta
	}

	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID, epsilon, delta)
	d := Sketch{info, sketch16, sync.RWMutex{}}
	err := d.Save()
	if err != nil {
		logger.Error.Printf("an error has occurred while saving Sketch: %s", err.Error())
	}
	return &d, nil
}

/*
NewSketchFromData ...
*/
func NewSketchFromData(info *abstract.Info) (*Sketch, error) {
	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID,
		info.Properties["epsilon"], info.Properties["delta"])
	// FIXME: create Sketch from new data
	return &Sketch{info, sketch16, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.IncreaseCount(value)
	d.Save()
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Sketch) AddMultiple(values [][]byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	for _, value := range values {
		d.impl.IncreaseCount(value)
	}
	d.Save()
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
	d.lock.RLock()
	defer d.lock.RUnlock()
	return uint(d.impl.Count())
}

/*
Clear ...
*/
func (d *Sketch) Clear() (bool, error) {
	d.impl.Reset()
	return true, nil
}

/*
Save ...
*/
func (d *Sketch) Save() error {
	count := d.impl.Count()
	d.Info.State["count"] = uint64(count)
	infoData, err := json.Marshal(d.Info)
	if err != nil {
		return err
	}
	return storage.GetManager().SaveInfo(d.Info.ID, infoData)
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
	res := make(map[string]uint)
	for _, value := range values {
		res[string(value)] = uint(d.impl.GetCount(value))
	}
	return res
}
