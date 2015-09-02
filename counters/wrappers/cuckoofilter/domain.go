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
Domain is the toplevel domain to control the HLL implementation
*/
type Domain struct {
	*abstract.Info
	impl *cuckoofilter.CuckooFilter
	lock sync.RWMutex
}

/*
NewDomain ...
*/
func NewDomain(info *abstract.Info) (*Domain, error) {
	d := &Domain{info, cuckoofilter.NewCuckooFilter(info), sync.RWMutex{}}
	d.Save()
	return d, nil
}

/*
Add ...
*/
func (d *Domain) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := d.impl.InsertUnique(value)
	err := d.Save()
	return ok, err
}

/*
AddMultiple ...
*/
func (d *Domain) AddMultiple(values [][]byte) (bool, error) {
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
func (d *Domain) Remove(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ok := d.impl.Delete(value)
	err := d.Save()
	return ok, err
}

/*
RemoveMultiple ...
*/
func (d *Domain) RemoveMultiple(values [][]byte) (bool, error) {
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
func (d *Domain) GetCount() uint {
	return uint(d.impl.GetCount())
}

/*
Clear ...
*/
func (d *Domain) Clear() (bool, error) {
	return true, nil
}

/*
Save ...
*/
func (d *Domain) Save() error {
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
func (d *Domain) GetType() string {
	return d.Type
}

/*
GetFrequency ...
*/
func (d *Domain) GetFrequency(values [][]byte) interface{} {
	return nil
}
