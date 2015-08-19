package counters

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/cuckoofilter"
	"github.com/seiflotfy/skizze/counters/wrappers/hllpp"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

/*
ManagerStruct is responsible for manipulating the counters and syncing to disk
*/
type ManagerStruct struct {
	domains map[string]abstract.Counter
	info    map[string]abstract.Info
}

var manager *ManagerStruct
var logger = utils.GetLogger()

/*
CreateDomain ...
*/
func (m *ManagerStruct) CreateDomain(domainID string, domainType string, capacity uint64) error {
	//TODO: spit errir uf domainType is invalid
	//FIXME: no hardcoding of immutable here
	if len([]byte(domainID)) > config.MaxKeySize {
		return errors.New("invalid length of domain ID: " + strconv.Itoa(len(domainID)) + ". Max length allowed: " + strconv.Itoa(config.MaxKeySize))
	}
	info := &abstract.Info{ID: domainID,
		Type:     domainType,
		Capacity: capacity,
		State:    make(map[string]uint64)}
	var domain abstract.Counter
	var err error
	switch domainType {
	case abstract.Cardinality:
		domain, err = hllpp.NewDomain(info)
	case abstract.PurgableCardinality:
		domain, err = cuckoofilter.NewDomain(info)
	default:
		return errors.New("invalid domain type: " + domainType)
	}

	if err != nil {
		errTxt := fmt.Sprint("Could not load domain ", info, ". Err:", err)
		return errors.New(errTxt)
	}
	m.domains[info.ID] = domain
	m.dumpInfo(info)
	return nil
}

/*
DeleteDomain ...
*/
func (m *ManagerStruct) DeleteDomain(domainID string) error {
	if _, ok := m.domains[domainID]; !ok {
		return errors.New("No such domain " + domainID)
	}
	delete(m.domains, domainID)
	//FIXME: delete from storage
	return nil
}

/*
GetDomains ...
*/
func (m *ManagerStruct) GetDomains() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	domains := make([]string, len(m.domains), len(m.domains))
	i := 0
	for k := range m.domains {
		domains[i] = k
		i++
	}
	return domains, nil
}

/*
AddToDomain ...
*/
func (m *ManagerStruct) AddToDomain(domainID string, values []string) error {
	var val, ok = m.domains[domainID]
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
	var val, ok = m.domains[domainID]
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
	var val, ok = m.domains[domainID]
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
	domains := make(map[string]abstract.Counter)
	manager = &ManagerStruct{domains, make(map[string]abstract.Info)}
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
		var domain abstract.Counter
		var err error
		switch info.Type {
		case abstract.Cardinality:
			domain, err = hllpp.NewDomainFromData(&info)
		case abstract.PurgableCardinality:
			domain, err = cuckoofilter.NewDomain(&info)
		default:
			logger.Info.Println("Invalid counter type", info.Type)
		}
		if err != nil {
			errTxt := fmt.Sprint("Could not load domain ", info, ". Err:", err)
			return errors.New(errTxt)
		}
		m.domains[info.ID] = domain
		strg.LoadData(key, 0, 0)
	}
	return nil
}
