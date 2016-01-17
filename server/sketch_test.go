package server

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/gogo/protobuf/proto"
	"github.com/skizzehq/skizze/config"
	pb "github.com/skizzehq/skizze/datamodel"
	"github.com/skizzehq/skizze/utils"
)

func TestCreateSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}
}

func TestCreateAddInvalidSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
		Defaults: &pb.Defaults{
			Capacity: proto.Int64(1337), // FIXME: Allow default as -1
			Rank:     proto.Int64(7),
		},
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	typ = pb.SketchType_FREQ
	in.Type = &typ
	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d", "a", "b"},
	}

	if _, err := client.Add(context.Background(), addReq); err == nil {
		t.Error("Expect error, got", err)
	}

}

func TestCreateAddDeleteAddSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d", "a", "b"},
	}

	if _, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	}

	typ = pb.SketchType_RANK
	if _, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	if res, err := client.ListAll(context.Background(), &pb.Empty{}); err != nil {
		t.Error("Did not expect error, got", err)
	} else if len(res.GetSketches()) != 2 {
		t.Error("Expected len(res) == 2, got ", len(res.GetSketches()))
	}

	if _, err := client.DeleteSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	}

	if _, err := client.Add(context.Background(), addReq); err == nil {
		t.Error("Expected error, got", err)
	}

	if res, err := client.ListAll(context.Background(), &pb.Empty{}); err != nil {
		t.Error("Did not expect error, got", err)
	} else if len(res.GetSketches()) != 1 {
		t.Error("Expected len(res) == 1, got ", len(res.GetSketches()))
	}

}

func TestAddGetCardSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_CARD
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d", "a", "b"},
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	getReq := &pb.GetRequest{
		Sketch: in,
		Values: []string{},
	}

	if res, err := client.GetCardinality(context.Background(), getReq); err != nil {
		t.Error("Did not expect error, got", err)

	} else if res.GetCardinality() != 4 {
		t.Error("Expected cardinality 4, got", res.GetCardinality())
	}
}

func TestAddGetMembSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_MEMB
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "a", "b", "c", "d"},
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	getReq := &pb.GetRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d", "e", "b"},
	}

	expected := map[string]bool{
		"a": true, "b": true, "c": true, "d": true, "e": false,
	}

	if res, err := client.GetMembership(context.Background(), getReq); err != nil {
		t.Error("Did not expect error, got", err)
	} else {
		for _, v := range res.Memberships {
			if expected[v.GetValue()] != v.GetIsMember() {
				t.Errorf("Expected %s == %t, got", v.GetValue(), v.GetIsMember())
			}
		}
	}
}

func TestAddGetFreqSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_FREQ
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "a", "b", "c", "d"},
	}

	expected := map[string]int64{
		"a": 2, "b": 1, "c": 1, "d": 1, "e": 0,
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	getReq := &pb.GetRequest{
		Sketch: in,
		Values: []string{"a", "b", "c", "d", "e", "b"},
	}

	if res, err := client.GetFrequency(context.Background(), getReq); err != nil {
		t.Error("Did not expect error, got", err)
	} else {
		for _, v := range res.GetFrequencies() {
			if expected[v.GetValue()] != v.GetCount() {
				t.Errorf("Expected %s == %d, got", v.GetValue(), v.GetCount())
			}
		}
	}
}

func TestAddGetRankSketch(t *testing.T) {
	config.Reset()
	utils.SetupTests()
	defer utils.TearDownTests()

	client, conn := setupClient()
	defer tearDownClient(conn)

	typ := pb.SketchType_RANK
	name := "yoyo"

	in := &pb.Sketch{
		Name: proto.String(name),
		Type: &typ,
	}

	if res, err := client.CreateSketch(context.Background(), in); err != nil {
		t.Error("Did not expect error, got", err)
	} else if res.GetName() != in.GetName() {
		t.Errorf("Expected name == %s, got %s", in.GetName(), res.GetName())
	} else if res.GetType() != in.GetType() {
		t.Errorf("Expected name == %q, got %q", in.GetType(), res.GetType())
	}

	addReq := &pb.AddRequest{
		Sketch: in,
		Values: []string{"a", "a", "b", "c", "d", "a", "b", "a", "b", "c"},
	}

	expected := map[string]int64{
		"a": 4, "b": 3, "c": 2, "d": 1, "e": 0,
	}

	if _, err := client.Add(context.Background(), addReq); err != nil {
		t.Error("Did not expect error, got", err)
	}

	getReq := &pb.GetRequest{
		Sketch: in,
	}

	if res, err := client.GetRank(context.Background(), getReq); err != nil {
		t.Error("Did not expect error, got", err)
	} else {
		for _, v := range res.GetRanks() {
			if expected[v.GetValue()] != v.GetCount() {
				t.Errorf("Expected %s == %d, got", v.GetValue(), v.GetCount())
			}
		}
	}
}
