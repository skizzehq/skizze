package hllpp

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/hllpp/hllpp"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

/*
Sketch is the toplevel sketch to control the HLL implementation
*/
type Sketch struct {
	*abstract.Info
	impl *hllpp.HLLPP
	lock sync.RWMutex
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	d := Sketch{info, hllpp.New(), sync.RWMutex{}}
	err := d.Save()
	if err != nil {
		logger.Error.Println("an error has occurred while saving sketch: " + err.Error())
	}
	return &d, nil
}

/*
NewSketchFromData ...
*/
func NewSketchFromData(info *abstract.Info) (*Sketch, error) {
	manager = storage.GetManager()
	data, err := manager.LoadData(info.ID, 0, 0)
	counter, err := hllpp.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &Sketch{info, counter, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.Add(value)
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
		d.impl.Add(value)
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
	return true, nil
}

/*
Save ...
*/
func (d *Sketch) Save() error {
	serialized := d.impl.Marshal()
	err := manager.SaveData(d.Info.ID, serialized, 0)
	if err != nil {
		return err
	}
	info, _ := json.Marshal(d.Info)
	return manager.SaveInfo(d.Info.ID, info)
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
	return nil
}
