package cuckoofilter

import (
	"encoding/json"

	"github.com/seiflotfy/counts/counters/abstract"
	"github.com/seiflotfy/counts/counters/wrappers/cuckoofilter/cuckoofilter"
	"github.com/seiflotfy/counts/storage"
	"github.com/seiflotfy/counts/utils"
)

var logger = utils.GetLogger()

/*
Domain is the toplevel domain to control the HLL implementation
*/
type Domain struct {
	abstract.Info
	impl *cuckoofilter.CuckooFilter
}

/*
NewDomain ...
*/
func NewDomain(info abstract.Info) (Domain, error) {
	return Domain{info, cuckoofilter.NewCuckooFilter(info)}, nil
}

/*
Add ...
*/
func (d Domain) Add(value []byte) (bool, error) {
	ok := d.impl.InsertUnique(value)
	err := d.Save()
	return ok, err
}

/*
AddMultiple ...
*/
func (d Domain) AddMultiple(values [][]byte) (bool, error) {
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
func (d Domain) Remove(value []byte) (bool, error) {
	ok := d.impl.Delete(value)
	err := d.Save()
	return ok, err
}

/*
RemoveMultiple ...
*/
func (d Domain) RemoveMultiple(values [][]byte) (bool, error) {
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
func (d Domain) GetCount() uint {
	return uint(d.impl.GetCount())
}

/*
Clear ...
*/
func (d Domain) Clear() (bool, error) {
	return true, nil
}

/*
Save ...
*/
func (d Domain) Save() error {
	count := d.impl.GetCount()
	d.Info.State["count"] = uint64(count)
	infoData, err := json.Marshal(d.Info)
	if err != nil {
		return err
	}
	err = storage.GetManager().SaveInfo(d.Info.ID, infoData)
	return err
}
