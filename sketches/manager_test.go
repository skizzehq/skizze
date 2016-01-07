package sketches

import (
	"path/filepath"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

func TestNewManager(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()
	if _, err := NewManager(); err != nil {
		t.Error("Expected no errors, got", err)
	}
}

func TestNoSketches(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()
	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestCreateSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

}

func TestCreateGetSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if err := m.AddToSketch(info, []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if res, err := m.GetFromSketch(info, nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(uint) != 4 {
		t.Error("Expected res = 4, got", res)
	}
}

func TestCreateDuplicateSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if err := m.CreateSketch(info); err == nil {
		t.Error("Expected error (duplicate sketch), got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}
}

func TestCreateInvalidSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = "N/A"
	if err := m.CreateSketch(info); err == nil {
		t.Error("Expected error invalid sketch, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestDeleteNonExistingSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.DeleteSketch(info); err == nil {
		t.Error("Expected errors deleting non-existing sketch, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestDeleteSketch(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}
	if err := m.DeleteSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestSaveLoad(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	m1, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m1.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m1.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.Save(map[string]*datamodel.Info{info.ID(): info}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m1.GetFromSketch(info, nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(uint) != 5 {
		t.Error("Expected res = 5, got", res)
	}

	path := filepath.Join(config.GetConfig().DataDir, "marvel.card")
	if exists, err := utils.Exists(path); err != nil {
		t.Error("Expected no errors, got", err)
	} else if !exists {
		t.Errorf("Expected file dump %s to exists, but apparently it doesn't", path)
	}

	m1.Destroy()

	m1, err = NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m1.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card], got", sketches)
	}

	if res, err := m1.GetFromSketch(info, nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(uint) != 4 {
		t.Error("Expected res = 4, got", res)
	}
}
