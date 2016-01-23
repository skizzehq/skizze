package storage

import (
	"fmt"
	"reflect"
	"testing"

	"config"
	"datamodel"
	"utils"
)

func TestSaveLoadInfo(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	infos := map[string]*datamodel.Info{}
	for i := 0; i < 10; i++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.MaxUniqueItems = 10000
		info.Name = fmt.Sprintf("marvel-%d", i)
		info.Type = datamodel.HLLPP
		infos[info.ID()] = info
	}

	if err := m.SaveInfo(infos); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()

	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if infoMap, err := m.LoadAllInfo(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(infoMap) != 10 {
		t.Error("Expected 10 info in map, go", len(infoMap))
	}
}

func TestSaveLoadDomain(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	domains := make(map[string][]string)
	for i := 0; i < 10; i++ {
		domains[fmt.Sprintf("test-%d", i%2)] = append(domains[fmt.Sprintf("test-%d", i%2)],
			fmt.Sprintf("test-%d", i))
	}

	if err := m.SaveDomains(domains); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()
	// Create new Manager
	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if doms, err := m.LoadAllDomains(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(domains) > 2 {
		t.Error("Expected 2 domains, got", len(doms))
	} else if eq := reflect.DeepEqual(domains, doms); !eq {
		t.Error("Expected domains == doms, got", eq)
	}
}

func TestOverwriteInfo(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	infos := map[string]*datamodel.Info{}
	for i := 0; i < 10; i++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.MaxUniqueItems = 10000
		info.Name = fmt.Sprintf("marvel-%d", i)
		info.Type = datamodel.HLLPP
		infos[info.ID()] = info
	}

	if err := m.SaveInfo(infos); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()

	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if infoMap, err := m.LoadAllInfo(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(infoMap) != 10 {
		t.Error("Expected 10 info in map, go", len(infoMap))
	}

	delete(infos, "marvel-0.card")
	if err := m.SaveInfo(infos); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()

	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if infoMap, err := m.LoadAllInfo(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(infoMap) != 9 {
		t.Error("Expected 9 info in map, go", len(infoMap))
	}
}

func TestOverwriteDomains(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	domains := make(map[string][]string)
	for i := 0; i < 10; i++ {
		domains[fmt.Sprintf("test-%d", i%2)] = append(domains[fmt.Sprintf("test-%d", i%2)],
			fmt.Sprintf("test-%d", i))
	}

	if err := m.SaveDomains(domains); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()
	// Create new Manager
	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if doms, err := m.LoadAllDomains(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(domains) > 2 {
		t.Error("Expected 2 domains, got", len(doms))
	} else if eq := reflect.DeepEqual(domains, doms); !eq {
		t.Error("Expected domains == doms, got", eq)
	}

	delete(domains, "test-0")

	if err := m.SaveDomains(domains); err != nil {
		t.Error("Expected no errors, got", err)
	}

	_ = m.Close()
	// Create new Manager
	m = NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}
	if doms, err := m.LoadAllDomains(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(domains) > 2 {
		t.Error("Expected 1 domain, got", len(doms))
	} else if eq := reflect.DeepEqual(domains, doms); !eq {
		t.Error("Expected domains == doms, got", eq)
	}
}

func TestEmptyLoadAllInfo(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	if info, err := m.LoadAllInfo(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(info) > 0 {
		t.Error("Expected no info, got", info)
	}
}

func TestEmptyLoadAllDomains(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	if domains, err := m.LoadAllDomains(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(domains) > 0 {
		t.Error("Expected no info, got", domains)
	}
}
