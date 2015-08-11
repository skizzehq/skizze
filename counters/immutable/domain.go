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
	d := Domain{info, hllpp.New()}
	d.Save()
	return d
}

/*
NewDomainFromData ...
*/
func NewDomainFromData(info abstract.Info) Domain {
	data := storage.GetManager().LoadData(info.ID, 0, 0)
	counter, err := hllpp.Unmarshal(data)
	utils.PanicOnError(err)
	return Domain{info, counter}
}

/*
Add ...
*/
func (d Domain) Add(value []byte) (bool, error) {
	d.impl.Add(value)
	d.Save()
	return true, nil
}

/*
AddMultiple ...
*/
func (d Domain) AddMultiple(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	d.Save()
	return true, nil
}

/*
Remove ...
*/
func (d Domain) Remove(value []byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
}

/*
RemoveMultiple ...
*/
func (d Domain) RemoveMultiple(values [][]byte) (bool, error) {
	logger.Error.Println("This domain type does not support deletion")
	return false, errors.New("This domain type does not support deletion")
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
