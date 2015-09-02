package cml

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/count-min-log/count-min-log"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

const defaultEpsilon = 0.00000543657
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
	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID, defaultEpsilon, defaultDelta)
	d := Sketch{info, sketch16, sync.RWMutex{}}
	err := d.Save()
	if err != nil {
		logger.Error.Println("an error has occurred while saving Sketch: " + err.Error())
	}
	return &d, nil
}

/*
NewSketchFromData ...
*/
func NewSketchFromData(info *abstract.Info) (*Sketch, error) {
	sketch16, _ := cml.NewSketch16ForEpsilonDelta(info.ID, defaultEpsilon, defaultDelta)
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
GetFrequency ...
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	res := make(map[string]uint)
	for _, value := range values {
		res[string(value)] = uint(d.impl.GetCount(value))
	}
	return res
}
