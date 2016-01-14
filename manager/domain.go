package manager

import (
	"fmt"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

type domainManager struct {
	domains  map[string][]string
	sketches *sketchManager
	info     *infoManager
	storage  *storage.Manager
}

func newDomainManager(info *infoManager, sketches *sketchManager, storage *storage.Manager) *domainManager {
	domains, err := storage.LoadAllDomains()
	utils.PanicOnError(err)
	return &domainManager{
		domains:  domains,
		info:     info,
		sketches: sketches,
		storage:  storage,
	}
}

func (m *domainManager) create(id string, infos map[string]*datamodel.Info) error {
	if _, ok := m.domains[id]; ok {
		return fmt.Errorf(`Domain with name "%s" already exists`, id)
	}

	var err error
	ids := make([]string, len(infos), len(infos))
	tmpInfos := make(map[string]*datamodel.Info)
	tmpSketches := make(map[string]*datamodel.Info)
	for id, info := range infos {
		if err = m.info.create(info); err != nil {
			break
		}
		tmpInfos[id] = info
		if err = m.sketches.create(info); err != nil {
			break
		}
		tmpSketches[id] = info

	}

	if len(tmpInfos) != len(infos) {
		for _, v := range tmpInfos {
			if err := m.info.delete(v.ID()); err != nil {
				// TODO: print out something
			}
		}
	}
	if len(tmpSketches) != len(infos) {
		for _, v := range tmpSketches {
			if err := m.sketches.delete(v.ID()); err != nil {
				// TODO: print out something
			}
		}
	}

	m.domains[id] = ids
	return err
}

func (m *domainManager) delete(id string) error {
	if ids, ok := m.domains[id]; ok {
		for _, id := range ids {
			if info := m.info.get(id); info != nil {
				if err := m.sketches.delete(info.ID()); err != nil {
					// TODO: print something ?
				}
				if err := m.info.delete(info.ID()); err != nil {
					// TODO: print something ?
				}
			}
		}
	}
	// FIXME: return error if not exist ?
	return nil
}

func (m *domainManager) save() error {
	return m.storage.SaveDomains(m.domains)
}
