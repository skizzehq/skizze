package sketches

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/count-min-log"
	"github.com/seiflotfy/skizze/sketches/wrappers/hllpp"
	"github.com/seiflotfy/skizze/sketches/wrappers/topk"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

/*
ManagerStruct is responsible for manipulating the sketches and syncing to disk
*/
type ManagerStruct struct {
	sketches map[string]abstract.Sketch
	info     map[string]*abstract.Info
}

var manager *ManagerStruct
var logger = utils.GetLogger()

/*
CreateSketch ...
*/
func (m *ManagerStruct) CreateSketch(sketchID string, sketchType string, props map[string]float64) error {
	id := fmt.Sprintf("%s.%s", sketchID, sketchType)

	// Check if sketch with ID already exists
	if info, ok := m.info[id]; ok {
		errStr := fmt.Sprintf("Sketch %s of type %s already exists", sketchID, info.Type)
		return errors.New(errStr)
	}

	// Check that id length does not exceed MaxKeySize
	if len([]byte(id)) > config.MaxKeySize {
		errStr := fmt.Sprintf("Invalid length of sketch ID: %d. Max length allowed: %d", len(id), config.MaxKeySize)
		return errors.New(errStr)
	}

	// Make sure sketchType is set
	if sketchType == "" {
		logger.Error.Println("SketchType is mandatory and must be set!")
		return errors.New("No sketch type was given!")
	}

	info := &abstract.Info{ID: id,
		Type:       sketchType,
		Properties: props,
		State:      make(map[string]uint64)}
	var sketch abstract.Sketch
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
	m.sketches[id] = sketch
	m.dumpInfo(info)
	return nil
}

/*
DeleteSketch ...
*/
func (m *ManagerStruct) DeleteSketch(sketchID string, sketchType string) error {
	id := fmt.Sprintf("%s.%s", sketchID, sketchType)

	if _, ok := m.sketches[id]; !ok {
		return errors.New("No such sketch " + sketchID)
	}
	delete(m.sketches, id)
	delete(m.info, id)
	manager := storage.GetManager()
	err := manager.DeleteInfo(id)
	if err != nil {
		return err
	}
	return manager.DeleteData(id)
}

/*
GetSketches ...
*/
func (m *ManagerStruct) GetSketches() ([]string, error) {
	// TODO: Remove dummy data and implement proper result
	sketches := make([]string, len(m.sketches), len(m.sketches))
	i := 0
	for _, v := range m.sketches {
		typ := v.GetType()
		id := v.GetID()
		sketches[i] = fmt.Sprintf("%s/%s", typ, id[:len(id)-len(typ)-1])
		i++
	}
	return sketches, nil
}

/*
AddToSketch ...
*/
func (m *ManagerStruct) AddToSketch(sketchID string, sketchType string, values []string) error {
	id := fmt.Sprintf("%s.%s", sketchID, sketchType)

	var val, ok = m.sketches[id]
	if ok == false {
		errStr := fmt.Sprintf("No such sketch %s of type %s found", sketchID, sketchType)
		return errors.New(errStr)
	}
	var sketch abstract.Sketch
	sketch = val.(abstract.Sketch)

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	_, err := sketch.AddMultiple(bytes)
	return err
}

/*
DeleteFromSketch ...
*/
func (m *ManagerStruct) DeleteFromSketch(sketchID string, sketchType string, values []string) error {
	var val, ok = m.sketches[sketchID]
	if ok == false {
		return errors.New("No such sketch: " + sketchID)
	}
	var sketch abstract.Sketch
	sketch = val.(abstract.Sketch)

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	ok, err := sketch.RemoveMultiple(bytes)
	return err
}

/*
GetCountForSketch ...
*/
func (m *ManagerStruct) GetCountForSketch(sketchID string, sketchType string, values []string) (interface{}, error) {
	id := fmt.Sprintf("%s.%s", sketchID, sketchType)
	var val, ok = m.sketches[id]
	if ok == false {
		errStr := fmt.Sprintf("No such sketch %s of type %s found", sketchID, sketchType)
		return 0, errors.New(errStr)
	}
	var sketch abstract.Sketch
	sketch = val.(abstract.Sketch)

	if sketch.GetType() == abstract.CML {
		bvalues := make([][]byte, len(values), len(values))
		for i, value := range values {
			bvalues[i] = []byte(value)
		}
		return sketch.GetFrequency(bvalues), nil
	} else if sketch.GetType() == abstract.TopK {
		return sketch.GetFrequency(nil), nil
	}

	count := sketch.GetCount()
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
	sketches := make(map[string]abstract.Sketch)
	m := &ManagerStruct{sketches, make(map[string]*abstract.Info)}
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

func (m *ManagerStruct) dumpInfo(info *abstract.Info) {
	m.info[info.ID] = info
	manager := storage.GetManager()
	infoData, err := json.Marshal(info)
	utils.PanicOnError(err)
	manager.SaveInfo(info.ID, infoData)
}

func (m *ManagerStruct) loadInfo() error {
	manager := storage.GetManager()
	infos, err := manager.LoadAllInfo()
	if err != nil {
		return err
	}
	for _, infoData := range infos {
		var infoStruct abstract.Info
		err := json.Unmarshal(infoData, &infoStruct)
		if err != nil {
			return err
		}
		m.info[infoStruct.ID] = &infoStruct
	}
	return nil
}

func (m *ManagerStruct) loadSketches() error {
	for _, info := range m.info {
		var sketch abstract.Sketch
		var err error
		switch info.Type {
		case abstract.HLLPP:
			sketch, err = hllpp.NewSketchFromData(info)
		case abstract.TopK:
			sketch, err = topk.NewSketchFromData(info)
		case abstract.CML:
			sketch, err = cml.NewSketchFromData(info)
		default:
			logger.Info.Println("Invalid sketch type", info.Type)
		}
		if err != nil {
			errTxt := fmt.Sprint("Could not load sketch ", info, ". Err: ", err)
			return errors.New(errTxt)
		}
		m.sketches[info.ID] = sketch
	}
	return nil
}

/*
Destroy ...
*/
func (m *ManagerStruct) Destroy() {
	manager = nil
}
