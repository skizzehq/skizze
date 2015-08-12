package mutable

import (
	"github.com/seiflotfy/counts/counters/abstract"
	"github.com/seiflotfy/counts/counters/mutable/cuckoofilter"
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
func NewDomain(info abstract.Info) Domain {
	return Domain{info, cuckoofilter.NewCuckooFilter(uint(info.Capacity))}
}

/*
Add ...
*/
func (d Domain) Add(value []byte) (bool, error) {
	ok := d.impl.InsertUnique(value)
	return ok, nil
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
	return ok, nil
}

/*
Remove ...
*/
func (d Domain) Remove(value []byte) (bool, error) {
	ok := d.impl.Delete(value)
	return ok, nil
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
	return ok, nil
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
