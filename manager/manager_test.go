package manager

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/skizzehq/skizze/config"
	"github.com/skizzehq/skizze/datamodel"
	"github.com/skizzehq/skizze/utils"
)

func TestNoSketches(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()
	m := NewManager()
	if sketches := m.GetSketches(); len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestCreateSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
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

	// Create a second Sketch
	info = datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.TopK
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 2 {
		t.Error("Expected 2 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches[0][0], sketches[0][1])
	} else if sketches[1][0] != "marvel" || sketches[1][1] != "rank" {
		t.Error("Expected [[marvel rank]], got", sketches[1][0], sketches[1][1])
	}
}

func TestCreateAndSaveSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()

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

	// Save state
	if err := m.Save(); err != nil {
		t.Error("Expected no errors, got", err)
	}

	// Create a second Sketch
	info = datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.TopK
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 2 {
		t.Error("Expected 2 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches[0][0], sketches[0][1])
	} else if sketches[1][0] != "marvel" || sketches[1][1] != "rank" {
		t.Error("Expected [[marvel rank]], got", sketches[1][0], sketches[1][1])
	}

	m.Destroy()

	// State should be equal before save
	m = NewManager()
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

}

func TestCreateDuplicateSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
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

	m := NewManager()
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

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.DeleteSketch(info.ID()); err == nil {
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

	m := NewManager()
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
	if err := m.DeleteSketch(info.ID()); err != nil {
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

	m := NewManager()

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.Save(); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
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

	m.Destroy()

	m = NewManager()
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card], got", sketches)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
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

	m := NewManager()

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.CML
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "freq" {
		t.Error("Expected [[marvel freq]], got", sketches)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.Save(); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
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

	m.Destroy()

	m = NewManager()
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "freq" {
		t.Error("Expected [[marvel freq], got", sketches)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(map[string]uint)["hulk"] != 1 {
		t.Error("Expected res = 1, got", res)
	}
}

func TestRankSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.TopK
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "rank" {
		t.Error("Expected [[marvel rank]], got", sketches)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.Save(); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
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

	m.Destroy()

	m = NewManager()
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "rank" {
		t.Error("Expected [[marvel rank], got", sketches)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
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

	m := NewManager()

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "marvel"
	info.Type = datamodel.Bloom
	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "memb" {
		t.Error("Expected [[marvel memb]], got", sketches)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.Save(); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "captian america", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.([]*datamodel.Member)) != 3 {
		t.Error("Expected len(res) = 3, got", len(res.(map[string]bool)))
	} else if v := res.([]*datamodel.Member)[0].Member; !v {
		t.Error("Expected 'hulk' == true , got", v)
	} else if v := res.([]*datamodel.Member)[1].Member; v {
		t.Error("Expected 'captian america' == false , got", v)
	} else if v := res.([]*datamodel.Member)[2].Member; !v {
		t.Error("Expected 'captian america' == true , got", v)
	}

	path := filepath.Join(config.GetConfig().DataDir, "marvel.memb")
	if exists, err := utils.Exists(path); err != nil {
		t.Error("Expected no errors, got", err)
	} else if !exists {
		t.Errorf("Expected file dump %s to exists, but apparently it doesn't", path)
	}

	m.Destroy()

	m = NewManager()
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketch, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "memb" {
		t.Error("Expected [[marvel memb], got", sketches)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "captian america", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.([]*datamodel.Member)) != 3 {
		t.Error("Expected len(res) = 3, got", len(res.(map[string]bool)))
	} else if v := res.([]*datamodel.Member)[0].Member; !v {
		t.Error("Expected 'hulk' == true , got", v)
	} else if v := res.([]*datamodel.Member)[1].Member; v {
		t.Error("Expected 'captian america' == false , got", v)
	} else if v := res.([]*datamodel.Member)[2].Member; v {
		t.Error("Expected 'captian america' == true , got", v)
	}
}

func TestCreateDomain(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 1000000
	info.Properties.Rank = 100
	info.Name = "marvel"
	if err := m.CreateDomain(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 4 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	// Create a second Sketch
	info = datamodel.NewEmptyInfo()
	info.Properties.Capacity = 10000
	info.Name = "dc"
	if err := m.CreateDomain(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 8 {
		t.Error("Expected 8 sketches, got", len(sketches))
	} else if sketches[0][0] != "dc" || sketches[0][1] != "card" {
		t.Error("Expected [[dc card]], got", sketches[0][0], sketches[0][1])
	} else if sketches[1][0] != "dc" || sketches[1][1] != "freq" {
		t.Error("Expected [[dc freq]], got", sketches[1][0], sketches[1][1])
	}
}
