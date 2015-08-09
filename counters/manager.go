package counters

import (
	"counts/counters/abstract"
	"counts/counters/immutable"
	"counts/counters/mutable"
	"errors"

	"github.com/hashicorp/golang-lru"
)

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
type managerStruct struct {
	cache *lru.Cache
}

var manager *managerStruct

/*
CreateDomain ...
*/
func (m *managerStruct) CreateDomain(domainID string, domainType string, capacity uint64) error {
	//TODO: spit errir uf domainType is invalid
	//FIXME: no hardcoding of immutable here

	info := abstract.Info{ID: domainID, Type: domainType, Capacity: capacity}
	switch domainType {
	case "immutable":
		m.cache.Add(info.ID, immutable.NewDomain(info))
	case "mutable":
		m.cache.Add(info.ID, mutable.NewDomain(info))
	default:
		return errors.New("invalid domain type: " + domainType)
	}
	return nil
}

/*
DeleteDomain ...
*/
func (m *managerStruct) DeleteDomain(domainID string) error {
	m.cache.Remove(domainID)
	return nil
}

/*
GetDomains ...
*/
func (m *managerStruct) GetDomains() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	values := manager.cache.Keys()
	domains := make([]string, len(values), len(values))
	for i, v := range values {
		domains[i] = v.(string)
	}
	return domains, nil
}

/*
AddToDomain ...
*/
func (m *managerStruct) AddToDomain(domainID string, values []string) error {
	var val, ok = m.cache.Get(domainID)
	if ok == false {
		return errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	counter.AddMultiple(bytes)
	return nil
}

/*
DeleteFromDomain ...
*/
func (m *managerStruct) DeleteFromDomain(domainID string, values []string) error {
	var val, ok = m.cache.Get(domainID)
	if ok == false {
		return errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	ok, err := counter.RemoveMultiple(bytes)
	return err
}

/*
GetCountForDomain ...
*/
func (m *managerStruct) GetCountForDomain(domainID string) (uint, error) {

	var val, ok = m.cache.Get(domainID)
	if ok == false {
		return 0, errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)
	count := counter.GetCount()
	return count, nil
}

/*
getManager returns a singleton Manager
*/
func getManager() *managerStruct {
	if manager == nil {
		cache, _ := lru.New(100)
		manager = &managerStruct{cache}
	}
	return manager
}

/*
Manager
*/
var Manager = getManager()
