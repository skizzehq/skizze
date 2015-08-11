package counters

import (
	"counts/counters/abstract"
	"counts/counters/immutable"
	"counts/counters/mutable"
	"counts/storage"
	"counts/utils"
	"encoding/json"
	"errors"

	"github.com/hashicorp/golang-lru"
)

/*
ManagerStruct is responsible for manipulating the counters and syncing to disk
*/
type ManagerStruct struct {
	cache *lru.Cache
	info  map[string]abstract.Info
}

var manager *ManagerStruct
var logger = utils.GetLogger()

/*
CreateDomain ...
*/
func (m *ManagerStruct) CreateDomain(domainID string, domainType string, capacity uint64) error {
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
	dumpInfo(&info)
	return nil
}

/*
DeleteDomain ...
*/
func (m *ManagerStruct) DeleteDomain(domainID string) error {
	m.cache.Remove(domainID)
	return nil
}

/*
GetDomains ...
*/
func (m *ManagerStruct) GetDomains() ([]string, error) {
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
func (m *ManagerStruct) AddToDomain(domainID string, values []string) error {
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
func (m *ManagerStruct) DeleteFromDomain(domainID string, values []string) error {
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
func (m *ManagerStruct) GetCountForDomain(domainID string) (uint, error) {
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
GetManager returns a singleton Manager
*/
func GetManager() *ManagerStruct {
	if manager == nil {
		manager = newManager()
	}
	return manager
}

func newManager() *ManagerStruct {
	cache, _ := lru.New(100)
	manager = &ManagerStruct{cache, make(map[string]abstract.Info)}
	manager.loadInfo()
	manager.loadData()
	return manager
}

func dumpInfo(i *abstract.Info) {
	manager := storage.GetManager()
	infoData, err := json.Marshal(i)
	utils.PanicOnError(err)
	manager.PutInfo(i.ID, infoData)
}

func (m *ManagerStruct) loadInfo() {
	manager := storage.GetManager()
	var infoStruct abstract.Info
	for _, infoData := range manager.GetAllInfo() {
		json.Unmarshal(infoData, &infoStruct)
		m.info[infoStruct.ID] = infoStruct
	}
}

func (m *ManagerStruct) loadData() {
	strg := storage.GetManager()
	for key, info := range m.info {
		switch info.Type {
		case "immutable":
			m.cache.Add(info.ID, immutable.NewDomainFromData(info))
		default:
			logger.Info.Println("Invalid counter type", info.Type)
		}
		strg.LoadData(key, 0, 0)
	}
}
