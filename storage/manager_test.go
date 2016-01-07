package storage

import (
	"fmt"
	"testing"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

func TestNewManager(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()
	if m := NewManager(); m == nil {
		t.Error("Expected m != nil, got", m)
	}
}

func TestSaveInfo(t *testing.T) {
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
