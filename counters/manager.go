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
	sketchs map[string]abstract.Counter
	info    map[string]*abstract.Info
}

var manager *ManagerStruct
var logger = utils.GetLogger()

/*
CreateSketch ...
*/
func (m *ManagerStruct) CreateSketch(sketchID string, sketchType string, capacity uint64) error {

	// Check if sketch with ID already exists
	if info, ok := m.info[sketchID]; ok {
		errStr := fmt.Sprintf("Sketch %s of type %s already exists", sketchID, info.Type)
		return errors.New(errStr)
	}

	if len([]byte(sketchID)) > config.MaxKeySize {
		return errors.New("Invalid length of sketch ID: " + strconv.Itoa(len(sketchID)) + ". Max length allowed: " + strconv.Itoa(config.MaxKeySize))
	}
	if sketchType == "" {
		logger.Error.Println("SketchType is mandatory and must be set!")
		return errors.New("No sketch type was given!")
	}
	info := &abstract.Info{ID: sketchID,
		Type:     sketchType,
		Capacity: capacity,
		State:    make(map[string]uint64)}
	var sketch abstract.Counter
	var err error
	switch sketchType {
	case abstract.HLLPP:
		sketch, err = hllpp.NewSketch(info)
	case abstract.TopK:
		sketch, err = topk.NewSketch(info)
	case abstract.CML:
		sketch, err = cml.NewSketch(info)
	default:
		return errors.New("Invalid sketch type: " + sketchType)
	}

	if err != nil {
		errTxt := fmt.Sprint("Could not load sketch ", info, ". Err:", err)
		return errors.New(errTxt)
	}
	m.sketchs[info.ID] = sketch
	m.dumpInfo(info)
	return nil
}

/*
DeleteSketch ...
*/
func (m *ManagerStruct) DeleteSketch(sketchID string) error {
	if _, ok := m.sketchs[sketchID]; !ok {
		return errors.New("No such sketch " + sketchID)
	}
	delete(m.sketchs, sketchID)
	delete(m.info, sketchID)
	manager := storage.GetManager()
	err := manager.DeleteInfo(sketchID)
	if err != nil {
		return err
	}
	return manager.DeleteData(sketchID)
}

/*
GetSketches ...
*/
func (m *ManagerStruct) GetSketches() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	sketchs := make([]string, len(m.sketchs), len(m.sketchs))
	i := 0
	for k := range m.sketchs {
		sketchs[i] = k
		i++
	}
	return sketchs, nil
}

/*
AddToSketch ...
*/
func (m *ManagerStruct) AddToSketch(sketchID string, values []string) error {
	var val, ok = m.sketchs[sketchID]
	if ok == false {
		return errors.New("No such sketch: " + sketchID)
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
func (m *ManagerStruct) DeleteFromSketch(sketchID string, values []string) error {
	var val, ok = m.sketchs[sketchID]
	if ok == false {
		return errors.New("No such sketch: " + sketchID)
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
func (m *ManagerStruct) GetCountForSketch(sketchID string, values []string) (interface{}, error) {
	var val, ok = m.sketchs[sketchID]
	if ok == false {
		return 0, errors.New("No such sketch: " + sketchID)
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
	sketchs := make(map[string]abstract.Counter)
	m := &ManagerStruct{sketchs, make(map[string]*abstract.Info)}
	err := m.loadInfo()
	if err != nil {
		return nil, err
	}
	err = m.loadSketches()
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

func (m *ManagerStruct) loadSketches() error {
	strg := storage.GetManager()
	for key, info := range m.info {
		var sketch abstract.Counter
		var err error
		switch info.Type {
		case abstract.HLLPP:
			sketch, err = hllpp.NewSketchFromData(info)
		case abstract.TopK:
			sketch, err = topk.NewSketchFromData(info)
		case abstract.CML:
			sketch, err = cml.NewSketchFromData(info)
		default:
			logger.Info.Println("Invalid counter type", info.Type)
		}
		if err != nil {
			errTxt := fmt.Sprint("Could not load sketch ", info, ". Err: ", err)
			return errors.New(errTxt)
		}
		m.sketchs[info.ID] = sketch
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
