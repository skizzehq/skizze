package counters

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/count-min-log"
	"github.com/seiflotfy/skizze/counters/wrappers/hllpp"
	"github.com/seiflotfy/skizze/counters/wrappers/topk"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

/*
ManagerStruct is responsible for manipulating the counters and syncing to disk
*/
type ManagerStruct struct {
	domains map[string]abstract.Counter
	info    map[string]*abstract.Info
}

var manager *ManagerStruct
var logger = utils.GetLogger()

/*
CreateSketch ...
*/
func (m *ManagerStruct) CreateSketch(domainID string, domainType string, capacity uint64) error {

	// Check if domain with ID already exists
	if info, ok := m.info[domainID]; ok {
		errStr := fmt.Sprintf("Sketch %s of type %s already exists", domainID, info.Type)
		return errors.New(errStr)
	}

	if len([]byte(domainID)) > config.MaxKeySize {
		return errors.New("Invalid length of domain ID: " + strconv.Itoa(len(domainID)) + ". Max length allowed: " + strconv.Itoa(config.MaxKeySize))
	}
	if domainType == "" {
		logger.Error.Println("SketchType is mandatory and must be set!")
		return errors.New("No domain type was given!")
	}
	info := &abstract.Info{ID: domainID,
		Type:     domainType,
		Capacity: capacity,
		State:    make(map[string]uint64)}
	var domain abstract.Counter
	var err error
	switch domainType {
	case abstract.HLLPP:
		domain, err = hllpp.NewSketch(info)
	case abstract.TopK:
		domain, err = topk.NewSketch(info)
	case abstract.CML:
		domain, err = cml.NewSketch(info)
	default:
		return errors.New("Invalid domain type: " + domainType)
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
DeleteSketch ...
*/
func (m *ManagerStruct) DeleteSketch(domainID string) error {
	if _, ok := m.domains[domainID]; !ok {
		return errors.New("No such domain " + domainID)
	}
	delete(m.domains, domainID)
	delete(m.info, domainID)
	manager := storage.GetManager()
	err := manager.DeleteInfo(domainID)
	if err != nil {
		return err
	}
	return manager.DeleteData(domainID)
}

/*
GetSketchs ...
*/
func (m *ManagerStruct) GetSketchs() ([]string, error) {
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
AddToSketch ...
*/
func (m *ManagerStruct) AddToSketch(domainID string, values []string) error {
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
DeleteFromSketch ...
*/
func (m *ManagerStruct) DeleteFromSketch(domainID string, values []string) error {
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
GetCountForSketch ...
*/
func (m *ManagerStruct) GetCountForSketch(domainID string, values []string) (interface{}, error) {
	var val, ok = m.domains[domainID]
	if ok == false {
		return 0, errors.New("No such domain: " + domainID)
	}
	var counter abstract.Counter
	counter = val.(abstract.Counter)

	if counter.GetType() == abstract.CML {
		bvalues := make([][]byte, len(values), len(values))
		for i, value := range values {
			bvalues[i] = []byte(value)
		}
		return counter.GetFrequency(bvalues), nil
	} else if counter.GetType() == abstract.TopK {
		return counter.GetFrequency(nil), nil
	}

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
	m := &ManagerStruct{domains, make(map[string]*abstract.Info)}
	err := m.loadInfo()
	if err != nil {
		return nil, err
	}
	err = m.loadSketchs()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ManagerStruct) dumpInfo(i *abstract.Info) {
	m.info[i.ID] = i
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
		m.info[infoStruct.ID] = &infoStruct
	}
	return nil
}

func (m *ManagerStruct) loadSketchs() error {
	strg := storage.GetManager()
	for key, info := range m.info {
		var domain abstract.Counter
		var err error
		switch info.Type {
		case abstract.HLLPP:
			domain, err = hllpp.NewSketchFromData(info)
		case abstract.TopK:
			domain, err = topk.NewSketchFromData(info)
		case abstract.CML:
			domain, err = cml.NewSketchFromData(info)
		default:
			logger.Info.Println("Invalid counter type", info.Type)
		}
		if err != nil {
			errTxt := fmt.Sprint("Could not load domain ", info, ". Err: ", err)
			return errors.New(errTxt)
		}
		m.domains[info.ID] = domain
		strg.LoadData(key, 0, 0)
	}
	return nil
}

/*
Destroy ...
*/
func (m *ManagerStruct) Destroy() {
	manager = nil
}
