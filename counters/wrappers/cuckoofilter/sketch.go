package cuckoofilter

import (
	"encoding/json"
	"sync"

	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/cuckoofilter/cuckoofilter"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

/*
Sketch is the toplevel sketch to control the HLL implementation
*/
type Sketch struct {
	*abstract.Info
	impl *cuckoofilter.CuckooFilter
	lock sync.RWMutex
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	d := &Sketch{info, cuckoofilter.NewCuckooFilter(info), sync.RWMutex{}}
	d.Save()
	return d, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := d.impl.InsertUnique(value)
	err := d.Save()
	return ok, err
}

/*
AddMultiple ...
*/
func (d *Sketch) AddMultiple(values [][]byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := true
	for _, value := range values {
		okk := d.impl.InsertUnique(value)
		if okk == false {
			ok = okk
		}
	}
	err := d.Save()
	return ok, err
}

/*
Remove ...
*/
func (d *Sketch) Remove(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := d.impl.Delete(value)
	err := d.Save()
	return ok, err
}

/*
RemoveMultiple ...
*/
func (d *Sketch) RemoveMultiple(values [][]byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := true
	for _, value := range values {
		okk := d.impl.Delete(value)
		if okk == false {
			ok = okk
		}
	}
	err := d.Save()
	return ok, err
}

/*
GetCount ...
*/
func (d *Sketch) GetCount() uint {
	return uint(d.impl.GetCount())
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
	count := d.impl.GetCount()
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
	return nil
}
