package sketches

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

/*
ManagerStruct is responsible for manipulating the sketches and syncing to disk
*/
type ManagerStruct struct {
	sketches map[string]*SketchProxy
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

	sketch, err := createSketch(info)
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
	manager := storage.Manager()
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
		typ := v.Type
		id := v.ID
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
	var sketch *SketchProxy
	sketch = val

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	_, err := sketch.Add(bytes)
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
	var sketch *SketchProxy
	sketch = val

	bytes := make([][]byte, len(values), len(values))
	for i, value := range values {
		bytes[i] = []byte(value)
	}
	ok, err := sketch.Remove(bytes)
	return err
}

/*
GetCountForSketch ...
*/
func (m *ManagerStruct) GetCountForSketch(sketchID string, sketchType string, values []string) (map[string]interface{}, error) {
	id := fmt.Sprintf("%s.%s", sketchID, sketchType)
	var val, ok = m.sketches[id]
	if ok == false {
		errStr := fmt.Sprintf("No such sketch %s of type %s found", sketchID, sketchType)
		return nil, errors.New(errStr)
	}
	var sketch *SketchProxy
	sketch = val
	count := sketch.Count(values)
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
	sketches := make(map[string]*SketchProxy)
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
	// FIXME: Should we panic here?
	m.info[info.ID] = info
	manager := storage.Manager()
	infoData, err := json.Marshal(info)
	utils.PanicOnError(err)
	err = manager.SaveInfo(info.ID, infoData)
	utils.PanicOnError(err)
}

func (m *ManagerStruct) loadInfo() error {
	manager := storage.Manager()
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
		sketch, err := loadSketch(info)
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
