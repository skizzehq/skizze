package manager

import (
	"fmt"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
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
			info.Type, info.Name)
	}
	m.info[info.ID()] = info
	return nil
}

// FIXME: should take array or map instead?
func (m *infoManager) delete(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; !ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already exists`,
			info.Type, info.Name)
	}
	// FIXME: return error if not exist
	delete(m.info, info.ID())
	return nil
}

func (m *infoManager) save(infos map[string]*datamodel.Info) error {
	if infos == nil || len(infos) == 0 {
		infos = m.info
	}
	return m.storage.SaveInfo(infos)
}

func (m *infoManager) getCopy() map[string]*datamodel.Info {
	return (map[string]*datamodel.Info)(m.info)
}
