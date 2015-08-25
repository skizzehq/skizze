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

/*
Domain is the toplevel domain to control the count-min-log implementation
*/
type Domain struct {
	*abstract.Info
	impl *cml.Sketch16
	lock sync.RWMutex
}

/*
NewDomain ...
*/
func NewDomain(info *abstract.Info) (*Domain, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	sketch16, _ := cml.NewDefaultSketch16()
	d := Domain{info, sketch16, sync.RWMutex{}}
	err := d.Save()
	if err != nil {
		logger.Error.Println("an error has occurred while saving domain: " + err.Error())
	}
	return &d, nil
}

/*
NewDomainFromData ...
*/
func NewDomainFromData(info *abstract.Info) (*Domain, error) {
	data, err := storage.GetManager().LoadData(info.ID, 0, 0)
	if err != nil {
		return nil, err
	}
	sketch16, _ := cml.NewDefaultSketch16()
	// FIXME: create domain from new data
	sketch16.GetCount(data)
	return &Domain{info, sketch16, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Domain) Add(value []byte) (bool, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.impl.IncreaseCount(value)
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
		d.impl.IncreaseCount(value)
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
func (d *Domain) GetCount() interface{} {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return uint(d.impl.Count())
}

/*
Clear ...
*/
func (d *Domain) Clear() (bool, error) {
	d.impl.Reset()
	return true, nil
}

/*
Save ...
*/
func (d *Domain) Save() error {
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
func (d *Domain) GetType() string {
	return d.Type
}
