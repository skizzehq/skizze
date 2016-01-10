package manager

import (
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

func TestCreate(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	s := storage.NewManager()
	m := newInfoManager(s)

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP

	if err := m.create(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

}

func TestSaveDelete(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	s := storage.NewManager()
	m := newInfoManager(s)

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP

	if err := m.create(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	// Save state
	if err := m.save(info.ID()); err != nil {
		t.Error("Expected no errors, got", err)
	}

	// delete old Info
	if err := m.delete(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

}

func TestCreateDuplicate(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	s := storage.NewManager()
	m := newInfoManager(s)

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP

	if err := m.create(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	// Save state
	if err := m.save(info.ID()); err != nil {
		t.Error("Expected no errors, got", err)
	}

	info2 := datamodel.NewEmptyInfo()
	info2.Properties.Capacity = 10000
	info2.Name = "marvel"
	info2.Type = datamodel.HLLPP

	if err := m.create(info2); err == nil {
		t.Error("Expected errors, got", err)
	}

}

func TestDeleteInvalid(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	s := storage.NewManager()
	m := newInfoManager(s)

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP

	if err := m.delete(info); err == nil {
		t.Error("Expected errors, got", err)
	}

}
