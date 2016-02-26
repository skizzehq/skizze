package manager

import (
	"fmt"
	"sort"
	"strconv"

	"datamodel"
	pb "datamodel/protobuf"

	"github.com/njpatel/loggo"
)

var logger = loggo.GetLogger("manager")

func isValidType(info *datamodel.Info) bool {
	if info.Type == nil {
		return false
	}
	return len(datamodel.GetTypeString(info.GetType())) != 0
}

// Manager is responsible for manipulating the sketches and syncing to disk
type Manager struct {
	infos    *infoManager
	sketches *sketchManager
	domains  *domainManager
}

// NewManager ...
func NewManager() *Manager {
	sketches := newSketchManager()
	infos := newInfoManager()
	domains := newDomainManager(infos, sketches)

	m := &Manager{
		sketches: sketches,
		infos:    infos,
		domains:  domains,
	}

	return m
}

// CreateSketch ...
func (m *Manager) CreateSketch(info *datamodel.Info) error {
	if !isValidType(info) {
		return fmt.Errorf("Can not create sketch of type %s, invalid type.", info.Type)
	}
	if err := m.infos.create(info); err != nil {
		return err
	}
	if err := m.sketches.create(info); err != nil {
		// If error occurred during creation of sketch, delete info
		if err2 := m.infos.delete(info.ID()); err2 != nil {
			return fmt.Errorf("%q\n%q ", err, err2)
		}
		return err
	}
	return nil
}

// CreateDomain ...
func (m *Manager) CreateDomain(info *datamodel.Info) error {
	infos := make(map[string]*datamodel.Info)
	for _, typ := range datamodel.GetTypesPb() {
		styp := typ
		tmpInfo := info.Copy()
		tmpInfo.Type = &styp
		infos[tmpInfo.ID()] = tmpInfo
	}
	return m.domains.create(info.GetName(), infos)
}

// AddToSketch ...
func (m *Manager) AddToSketch(id string, values []string) error {
	return m.sketches.add(id, values)
}

// AddToDomain ...
func (m *Manager) AddToDomain(id string, values []string) error {
	return m.domains.add(id, values)
}

// DeleteSketch ...
func (m *Manager) DeleteSketch(id string) error {
	if err := m.infos.delete(id); err != nil {
		return err
	}
	return m.sketches.delete(id)
}

// DeleteDomain ...
func (m *Manager) DeleteDomain(id string) error {
	return m.domains.delete(id)
}

type tupleResult [][2]string

func (slice tupleResult) Len() int {
	return len(slice)
}

func (slice tupleResult) Less(i, j int) bool {
	if slice[i][0] == slice[j][0] {
		return slice[i][1] < slice[j][1]
	}
	return slice[i][0] < slice[j][0]
}

func (slice tupleResult) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// GetSketches return a list of sketch tuples [name, type]
func (m *Manager) GetSketches() [][2]string {
	sketches := tupleResult{}
	for _, v := range m.infos.info {
		sketches = append(sketches,
			[2]string{v.GetName(),
				datamodel.GetTypeString(v.GetType())})
	}
	sort.Sort(sketches)
	return sketches
}

// GetDomains return a list of sketch tuples [name, type]
func (m *Manager) GetDomains() [][2]string {
	domains := tupleResult{}
	for k, v := range m.domains.domains {
		domains = append(domains, [2]string{k, strconv.Itoa(len(v))})
	}
	sort.Sort(domains)
	return domains
}

// GetSketch ...
func (m *Manager) GetSketch(id string) (*datamodel.Info, error) {
	info := m.infos.get(id)
	if info == nil {
		return nil, fmt.Errorf("No such sketch %s", id)
	}
	return info, nil
}

// GetDomain ...
func (m *Manager) GetDomain(id string) (*pb.Domain, error) {
	return m.domains.get(id)
}

// GetFromSketch ...
func (m *Manager) GetFromSketch(id string, data interface{}) (interface{}, error) {
	return m.sketches.get(id, data)
}

// Destroy ...
func (m *Manager) Destroy() {
}
