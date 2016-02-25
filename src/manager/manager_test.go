package manager

import (
	"fmt"
	"testing"

	"config"
	"datamodel"
	pb "datamodel/protobuf"
	"utils"
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
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	// Create a second Sketch
	info2 := datamodel.NewEmptyInfo()
	typ2 := pb.SketchType_RANK
	info2.Properties.Size = utils.Int64p(10)
	info2.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info2.Type = &typ2

	if err := m.CreateSketch(info2); err != nil {
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
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

	if err := m.CreateSketch(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	// Create a second Sketch
	info2 := datamodel.NewEmptyInfo()
	typ2 := pb.SketchType_RANK
	info2.Properties.MaxUniqueItems = utils.Int64p(10000)
	info2.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info2.Type = &typ2

	if err := m.CreateSketch(info2); err != nil {
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

func TestCreateDuplicateSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp("avengers")
	info.Type = nil
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
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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
	typ := pb.SketchType_CARD
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(*pb.CardinalityResult).GetCardinality() != 5 {
		t.Error("Expected res = 5, got", res.(*pb.CardinalityResult).GetCardinality())
	}
}

func TestFreqSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	typ := pb.SketchType_FREQ
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "thor", "iron man", "hawk-eye"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if res.(*pb.FrequencyResult).GetFrequencies()[0].GetCount() != 2 {
		t.Error("Expected res = 2, got", res.(*pb.FrequencyResult).GetFrequencies()[0].GetCount())
	}
}

func TestRankSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	typ := pb.SketchType_RANK
	info.Properties.Size = utils.Int64p(10)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), nil); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.(*pb.RankingsResult).GetRankings()) != 5 {
		t.Error("Expected len(res) = 5, got", len(res.(*pb.RankingsResult).GetRankings()))
	} else if res.(*pb.RankingsResult).GetRankings()[0].GetValue() != "black widow" {
		t.Error("Expected 'black widow', got", res.(*pb.RankingsResult).GetRankings()[0].GetValue())
	}
}

func TestMembershipSaveLoad(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	typ := pb.SketchType_MEMB
	info.Properties.MaxUniqueItems = utils.Int64p(1000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))
	info.Type = &typ

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

	if err := m.AddToSketch(info.ID(), []string{"hulk", "black widow", "black widow", "black widow", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	}

	if res, err := m.GetFromSketch(info.ID(), []string{"hulk", "captian america", "black widow"}); err != nil {
		t.Error("Expected no errors, got", err)
	} else if len(res.(*pb.MembershipResult).GetMemberships()) != 3 {
		t.Error("Expected len(res) = 3, got", len(res.(*pb.MembershipResult).GetMemberships()))
	} else if v := res.(*pb.MembershipResult).GetMemberships()[0].GetIsMember(); !v {
		t.Error("Expected 'hulk' == true , got", v)
	} else if v := res.(*pb.MembershipResult).GetMemberships()[1].GetIsMember(); v {
		t.Error("Expected 'captian america' == false , got", v)
	} else if v := res.(*pb.MembershipResult).GetMemberships()[2].GetIsMember(); !v {
		t.Error("Expected 'captian america' == true , got", v)
	}
}

func TestCreateDeleteDomain(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	m := NewManager()
	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(10000)
	info.Properties.Size = utils.Int64p(10000)
	info.Name = utils.Stringp(fmt.Sprintf("marvel"))

	if err := m.CreateDomain(info); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 4 {
		t.Error("Expected 1 sketches, got", len(sketches))
	} else if sketches[0][0] != "marvel" || sketches[0][1] != "card" {
		t.Error("Expected [[marvel card]], got", sketches)
	}

	// Create a second Sketch
	info2 := datamodel.NewEmptyInfo()
	info2.Properties.MaxUniqueItems = utils.Int64p(10000)
	info2.Name = utils.Stringp("dc")
	if err := m.CreateDomain(info2); err != nil {
		t.Error("Expected no errors, got", err)
	}
	if sketches := m.GetSketches(); len(sketches) != 8 {
		t.Error("Expected 8 sketches, got", len(sketches))
	} else if sketches[0][0] != "dc" || sketches[0][1] != "card" {
		t.Error("Expected [[dc card]], got", sketches[0][0], sketches[0][1])
	} else if sketches[1][0] != "dc" || sketches[1][1] != "freq" {
		t.Error("Expected [[dc freq]], got", sketches[1][0], sketches[1][1])
		fmt.Println(sketches)
	}
}
