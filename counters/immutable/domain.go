package immutable

import (
	"counts/counters/abstract"
	"counts/counters/immutable/hllpp"
	"counts/storage"
	"counts/utils"
	"errors"
)

var logger = utils.GetLogger()
var manager *storage.ManagerStruct

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
	manager = storage.GetManager()
	manager.Create(info.ID)
	return Domain{info, hllpp.New()}
}

/*
Add ...
*/
func (d Domain) Add(value []byte) (bool, error) {
	d.impl.Add(value)
	return true, nil
}

/*
AddMultiple ...
*/
func (d Domain) AddMultiple(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

/*
Remove ...
*/
func (d Domain) Remove(value []byte) (bool, error) {
	logger.Error.Println("This operation does not deletion of counters")
	return false, errors.New("This operation does not deletion of counters")
}

/*
RemoveMultiple ...
*/
func (d Domain) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This operation does not deletion of counters")
	return false, errors.New("This operation does not deletion of counters")
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
func (d Domain) Clear() (bool, error) {
	return true, nil
}

/*
Save ...
*/
func (d Domain) Save() {
	serialized := d.impl.Marshal()
	manager.SaveData(d.Info.ID, serialized, 0)
}
