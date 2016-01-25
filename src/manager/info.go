package manager

import (
	"fmt"

	"datamodel"
	"storage"
	"utils"
)

type infoManager struct {
	storage *storage.Manager
	info    map[string]*datamodel.Info
}

func newInfoManager(storage *storage.Manager) *infoManager {
	info, err := storage.LoadAllInfo()
	utils.PanicOnError(err)
	return &infoManager{
		info:    info,
		storage: storage,
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

func (m *infoManager) save() error {
	return m.storage.SaveInfo(m.info)
}

func (m *infoManager) getCopy() map[string]*datamodel.Info {
	return (map[string]*datamodel.Info)(m.info)
}
