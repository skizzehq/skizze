package topk

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"sync"

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
	lock sync.RWMutex
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
	manager.Create(info.ID)
	if info.Properties["capacity"] == 0 {
		info.Properties["capacity"] = defaultCapacity
	}
	d := Sketch{info, topk.New(int(info.Properties["capacity"])), sync.RWMutex{}}
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
	if err != nil {
		return nil, err
	}
	var network bytes.Buffer // Stand-in for a network connection
	network.Write(data)
	dec := gob.NewDecoder(&network) // Will read from network.

	var counter topk.Stream
	err = dec.Decode(&counter)
	if err != nil {
		return nil, err
	}
	return &Sketch{info, &counter, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	str := string(value)
	d.impl.Insert(str, 1)
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
		str := string(value)
		d.impl.Insert(str, 1)
	}
	d.Save()
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
Save ...
*/
func (d *Sketch) Save() error {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	// Encode (send) the value.
	err := enc.Encode(d.impl)
	err = manager.SaveData(d.Info.ID, network.Bytes(), 0)
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
GetID ...
*/
func (d *Sketch) GetID() string {
	return d.ID
}

/*
GetFrequency ...
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	keys := d.impl.Keys()
	result := make([]ResultElement, len(keys), len(keys))
	for i, k := range keys {
		result[i] = ResultElement(k)
	}
	return result
}
