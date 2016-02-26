package manager

import (
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"

	"datamodel"
	pb "datamodel/protobuf"
)

type domainManager struct {
	domains  map[string][]string
	sketches *sketchManager
	info     *infoManager
}

func newDomainManager(info *infoManager, sketches *sketchManager) *domainManager {
	return &domainManager{
		domains:  make(map[string][]string),
		info:     info,
		sketches: sketches,
	}
}

func (m *domainManager) create(id string, infos map[string]*datamodel.Info) error {
	if _, ok := m.domains[id]; ok {
		return fmt.Errorf(`Domain with name "%s" already exists`, id)
	}

	var err error
	var ids []string
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
		ids = append(ids, info.ID())
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

	if len(ids) < len(datamodel.GetTypes()) {
		return fmt.Errorf("Not enough sketches")
	}
	m.domains[id] = ids
	return nil
}

// FIXME: maybe return a list of errors?
func (m *domainManager) delete(id string) error {
	var lastErr error
	if ids, ok := m.domains[id]; ok {
		for _, id := range ids {
			if info := m.info.get(id); info != nil {
				if err := m.sketches.delete(info.ID()); err != nil {
					// TODO: print something ?
					lastErr = err
				}
				if err := m.info.delete(info.ID()); err != nil {
					// TODO: print something ?
					lastErr = err
				}
			}
		}
	}
	delete(m.domains, id)
	// FIXME: return error if not exist ?
	return lastErr
}

func (m *domainManager) add(id string, values []string) error {
	sketches, ok := m.domains[id]

	if !ok {
		return fmt.Errorf(`Domain "%s" does not exists`, id)
	}

	var wg sync.WaitGroup
	wg.Add(len(sketches))

	for _, sketch := range sketches {
		go func(sk string) {
			if err := m.sketches.add(sk, values); err != nil {
				logger.Errorf("%q\n", err)
			}
			wg.Done()
		}(sketch)
	}

	wg.Wait()
	return nil
}

// FIXME: return all sketches with domain
func (m *domainManager) get(id string) (*pb.Domain, error) {
	sketchIds, ok := m.domains[id]
	if !ok {
		return nil, fmt.Errorf("Could not find domain %s", id)
	}
	var sketches []*pb.Sketch
	for _, id := range sketchIds {
		s := m.info.get(id)
		if s != nil {
			sketches = append(sketches, s.Sketch)
		}
	}
	domain := &pb.Domain{
		Name:     proto.String(id),
		Sketches: sketches,
	}
	return domain, nil
}
