package sketches

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

func TestNewManager(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	if _, err := NewManager(); err != nil {
		t.Error("Expected no errors, got", err)
	}
}

func TestNoSketches(t *testing.T) {
	config.Reset()
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
	config.Reset()
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
	config.Reset()
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
	config.Reset()
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
	config.Reset()
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
	config.Reset()
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
	config.Reset()
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

func TestCardSaveLoad(t *testing.T) {
	config.Reset()
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
	time.Sleep(time.Second * 2)
}

func TestFreqSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m1, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.CML
	if err := m1.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m1.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "freq" {
		t.Error("Expected [[marvel freq]], got", sketches)
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

	if res, err := m1.GetFromSketch(info, []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(map[string]uint)["hulk"] != 2 {
		t.Error("Expected res = 2, got", res)
	}

	path := filepath.Join(config.GetConfig().DataDir, "marvel.freq")
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
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "freq" {
		t.Error("Expected [[marvel freq], got", sketches)
	}

	if res, err := m1.GetFromSketch(info, []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(map[string]uint)["hulk"] != 1 {
		t.Error("Expected res = 1, got", res)
	}
}

func TestRankSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m1, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.TopK
	if err := m1.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m1.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "rank" {
		t.Error("Expected [[marvel rank]], got", sketches)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.Save(map[string]*datamodel.Info{info.ID(): info}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m1.GetFromSketch(info, nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.([]*datamodel.Element)) != 5 {
		t.Error("Expected len(res) = 5, got", len(res.([]*datamodel.Element)))
	} else if res.([]*datamodel.Element)[0].Key != "black widow" {
		t.Error("Expected 'black widow', got", res.([]*datamodel.Element)[0].Key)
	}

	path := filepath.Join(config.GetConfig().DataDir, "marvel.rank")
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
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "rank" {
		t.Error("Expected [[marvel rank], got", sketches)
	}

	if res, err := m1.GetFromSketch(info, nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.([]*datamodel.Element)) != 4 {
		t.Error("Expected len(res) = 4, got", len(res.([]*datamodel.Element)))
	} else if res.([]*datamodel.Element)[0].Key != "hulk" {
		t.Error("Expected 'hulk', got", res.([]*datamodel.Element)[0].Key)
	}
}

func TestMembershipSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m1, err := NewManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.Bloom
	if err := m1.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m1.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "memb" {
		t.Error("Expected [[marvel memb]], got", sketches)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.Save(map[string]*datamodel.Info{info.ID(): info}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m1.AddToSketch(info, []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m1.GetFromSketch(info, []string{"hulk", "captian america", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.(map[string]bool)) != 3 {
		t.Error("Expected len(res) = 3, got", len(res.(map[string]bool)))
	} else if v, _ := res.(map[string]bool)["hulk"]; !v {
		t.Error("Expected 'hulk' == true , got", v)
	} else if v, _ := res.(map[string]bool)["captian america"]; v {
		t.Error("Expected 'captian america' == false , got", v)
	} else if v, _ := res.(map[string]bool)["black widow"]; !v {
		t.Error("Expected 'captian america' == true , got", v)
	}

	path := filepath.Join(config.GetConfig().DataDir, "marvel.memb")
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
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "memb" {
		t.Error("Expected [[marvel memb], got", sketches)
	}

	if res, err := m1.GetFromSketch(info, []string{"hulk", "captian america", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.(map[string]bool)) != 3 {
		t.Error("Expected len(res) = 3, got", len(res.(map[string]bool)))
	} else if v, _ := res.(map[string]bool)["hulk"]; !v {
		t.Error("Expected 'hulk' == true , got", v)
	} else if v, _ := res.(map[string]bool)["captian america"]; v {
		t.Error("Expected 'captian america' == false , got", v)
	} else if v, _ := res.(map[string]bool)["black widow"]; v {
		t.Error("Expected 'captian america' == false , got", v)
	}
}

// TODO: Add tests for having several sketches, saving half way,
// getting all sketches, then reloading and getting all sketches,
// Expected the first half of the sketches to be there
