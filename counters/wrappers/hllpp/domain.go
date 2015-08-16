package hllpp

import (
	"errors"
	"sync"

	"github.com/seiflotfy/counts/counters/abstract"
	"github.com/seiflotfy/counts/counters/wrappers/hllpp/hllpp"
	"github.com/seiflotfy/counts/storage"
	"github.com/seiflotfy/counts/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

/*
Domain is the toplevel domain to control the HLL implementation
*/
type Domain struct {
	abstract.Info
	impl *hllpp.HLLPP
	lock sync.RWMutex
}

/*
NewDomain ...
*/
func NewDomain(info abstract.Info) (*Domain, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	d := Domain{info, hllpp.New(), sync.RWMutex{}}
	d.Save()
	//FIXME: check for error
	return &d, nil
}

/*
NewDomainFromData ...
*/
func NewDomainFromData(info abstract.Info) (*Domain, error) {
	data, err := storage.GetManager().LoadData(info.ID, 0, 0)
	counter, err := hllpp.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return &Domain{info, counter, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Domain) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.Add(value)
	d.Save()
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Domain) AddMultiple(values [][]byte) (bool, error) {
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
func (d *Domain) Remove(value []byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
}

/*
RemoveMultiple ...
*/
func (d *Domain) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
}

/*
GetCount ...
*/
func (d *Domain) GetCount() uint {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return uint(d.impl.Count())
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
	serialized := d.impl.Marshal()
	err := manager.SaveData(d.Info.ID, serialized, 0)
	return err
}
