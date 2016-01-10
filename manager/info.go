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

func (m *infoManager) create(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already exists`,
			info.Type, info.Name)
	}
	m.info[info.ID()] = info
	return nil
}

func (m *infoManager) delete(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; !ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already exists`,
			info.Type, info.Name)
	}
	// FIXME: return error if not exist
	delete(m.info, info.ID())
	return nil
}

func (m *infoManager) save(id string) error {
	infos := make(map[string]*datamodel.Info)
	if info, ok := m.info[id]; ok {
		infos[id] = info
		if err := m.storage.SaveInfo(infos); err != nil {
			// FIXME: shoudl we panic here or handle gracefully
			return err
		}
	}
	return nil
}
