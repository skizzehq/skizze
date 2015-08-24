package topk

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"sync"

	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/topk/go-topk"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

/*
Domain is the toplevel domain to control the HLL implementation
*/
type Domain struct {
	*abstract.Info
	impl *topk.Stream
	lock sync.RWMutex
}

/*
ResultElement ...
*/
type ResultElement topk.Element

/*
NewDomain ...
*/
func NewDomain(info *abstract.Info) (*Domain, error) {
	manager = storage.GetManager()
	manager.Create(info.ID)
	d := Domain{info, topk.New(int(info.Capacity)), sync.RWMutex{}}
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
	return &Domain{info, &counter, sync.RWMutex{}}, nil
}

/*
Add ...
*/
func (d *Domain) Add(value []byte) (bool, error) {
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
func (d *Domain) AddMultiple(values [][]byte) (bool, error) {
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
	keys := d.impl.Keys()
	result := make([]ResultElement, len(keys), len(keys))
	for i, k := range keys {
		result[i] = ResultElement(k)
	}
	return result
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
func (d *Domain) GetType() string {
	return d.Type
}
