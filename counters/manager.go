package counters

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/seiflotfy/counts/counters/abstract"
	"github.com/seiflotfy/counts/counters/wrappers/cuckoofilter"
	"github.com/seiflotfy/counts/counters/wrappers/hllpp"
	"github.com/seiflotfy/counts/storage"
	"github.com/seiflotfy/counts/utils"

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
	case abstract.Default:
		m.cache.Add(info.ID, hllpp.NewDomain(info))
	case abstract.Purgable:
		m.cache.Add(info.ID, cuckoofilter.NewDomain(info))
	default:
		return errors.New("invalid domain type: " + domainType)
	}
	m.dumpInfo(&info)
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
func GetManager() (*ManagerStruct, error) {
	var err error
	if manager == nil {
		manager, err = newManager()
	}
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func newManager() (*ManagerStruct, error) {
	cache, _ := lru.New(100)
	manager = &ManagerStruct{cache, make(map[string]abstract.Info)}
	err := manager.loadInfo()
	if err != nil {
		return nil, err
	}
	err = manager.loadDomains()
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (m *ManagerStruct) dumpInfo(i *abstract.Info) {
	manager := storage.GetManager()
	infoData, err := json.Marshal(i)
	utils.PanicOnError(err)
	manager.SaveInfo(i.ID, infoData)
}

func (m *ManagerStruct) loadInfo() error {
	manager := storage.GetManager()
	var infoStruct abstract.Info
	infos, err := manager.LoadAllInfo()
	//FIXME: Should we panic?
	if err != nil {
		return err
	}
	for _, infoData := range infos {
		json.Unmarshal(infoData, &infoStruct)
		m.info[infoStruct.ID] = infoStruct
	}
	return nil
}

func (m *ManagerStruct) loadDomains() error {
	strg := storage.GetManager()
	for key, info := range m.info {
		switch info.Type {
		case abstract.Default:
			domain, err := hllpp.NewDomainFromData(info)
			if err != nil {
				errTxt := fmt.Sprint("Could not load domain ", info, ". Err:", err)
				return errors.New(errTxt)
			}
			m.cache.Add(info.ID, domain)
		default:
			logger.Info.Println("Invalid counter type", info.Type)
		}
		strg.LoadData(key, 0, 0)
	}
	return nil
}
