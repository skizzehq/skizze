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
	err := manager.Create(info.ID)
	if err != nil {
		return nil, err
	}
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

	sketch16, _ := cml.NewSketch16ForEpsilonDelta(epsilon, delta)
	d := Sketch{info, sketch16, sync.RWMutex{}}
	err = d.Save()
	if err != nil {
		logger.Error.Printf("an error has occurred while saving Sketch: %s", err.Error())
	}
	return &d, nil
}

/*
NewSketchFromData ...
*/
func NewSketchFromData(info *abstract.Info) (*Sketch, error) {
	manager = storage.GetManager()
	b, err := manager.LoadData(info.ID, 0, 0)
	if err != nil {
		return nil, err
	}
	sketch16, _ := cml.Unmarshall16(b)
	return &Sketch{info, sketch16, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.IncreaseCount(value)
	err := d.Save()
	if err != nil {
		logger.Error.Println(err)
	}
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
	err := d.Save()
	if err != nil {
		logger.Error.Println(err)
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
	d.impl.Reset()
	return true, nil
}

/*
Save ...
*/
func (d *Sketch) Save() error {
	data, err := d.impl.Marshall()
	if err != nil {
		return err
	}
	err = manager.SaveData(d.Info.ID, data, 0)
	if err != nil {
		return err
	}
	count := d.impl.TotalCount()
	d.Info.State["count"] = uint64(count)
	infoData, err := json.Marshal(d.Info)
	if err != nil {
		return err
	}
	err = storage.GetManager().SaveInfo(d.Info.ID, infoData)
	return err
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
		count := d.impl.Frequency(value)
		res[string(value)] = uint(count)
	}
	return res
}
