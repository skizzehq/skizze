package storage

import (
	"fmt"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

func TestSaveInfo(t *testing.T) {
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
		info.Properties.Capacity = 10000
		info.Name = fmt.Sprintf("marvel-%d", i)
		info.Type = datamodel.HLLPP
		infos[info.ID()] = info
	}

	if err := m.SaveInfo(infos); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if infoMap, err := m.LoadAllInfo(); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(infoMap) != 10 {
		t.Error("Expected 10 info in map, go", len(infoMap))
	}
}

func TestSaveDomain(t *testing.T) {
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
		info.Properties.Capacity = 10000
		info.Name = fmt.Sprintf("marvel-%d", i)
		info.Type = datamodel.HLLPP
		infos[info.ID()] = info
	}

	err := m.SaveInfo(infos)

	if err != nil {
		t.Error("Expected no errors, got", err)
	}
}

func TestLoadAllInfo(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	infos, err := m.LoadAllInfo()

	if err != nil && infos != nil {
		t.Error("Expected no errors, got", err)
	}
}

func TestLoadAllDomains(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if m == nil {
		t.Error("Expected m != nil, got", m)
	}

	domains, err := m.LoadAllDomains()

	if err != nil && domains != nil {
		t.Error("Expected no errors, got", err)
	}
}