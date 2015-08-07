package immutable

import (
	"counts/counters/abstract"
	"hllpp"
)

/*
Domain is the toplevel domain to control the HLL implementation
*/
type Domain struct {
	abstract.Info
	impl *hllpp.HLLPP
}

/*
NewDomain ...
*/
func NewDomain(info abstract.Info) Domain {
	return Domain{info, hllpp.New()}
}

/*
Add ...
*/
func (d Domain) Add(value []byte) bool {
	d.impl.Add(value)
	return true
}

/*
AddMultiple ...
*/
func (d Domain) AddMultiple(values [][]byte) bool {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true
}

/*
Remove ...
*/
func (d Domain) Remove(value []byte) bool {
	return true
}

/*
RemoveMultiple ...
*/
func (d Domain) RemoveMultiple(values [][]byte) bool {
	return true
}

/*
GetCount ...
*/
func (d Domain) GetCount() uint {
	return uint(d.impl.Count())
}

/*
Clear ...
*/
func (d Domain) Clear() bool {
	return true
}
