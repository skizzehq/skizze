package manager

import (
	"fmt"

	"datamodel"
)

type infoManager struct {
	info map[string]*datamodel.Info
}

func newInfoManager() *infoManager {
	return &infoManager{
		info: make(map[string]*datamodel.Info),
	}
}

func (m *infoManager) get(id string) *datamodel.Info {
	info, _ := m.info[id]
	return info
}

func (m *infoManager) create(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already exists`,
			info.GetType(), info.GetName())
	}
	m.info[info.ID()] = info
	return nil
}

// FIXME: should take array or map instead?
func (m *infoManager) delete(id string) error {
	if _, ok := m.info[id]; !ok {
		return fmt.Errorf(`Sketch "%s" already exists`, id)
	}
	// FIXME: return error if not exist
	delete(m.info, id)
	return nil
}
