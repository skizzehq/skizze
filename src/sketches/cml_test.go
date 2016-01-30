package sketches

import (
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

func TestAdd(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1000000)
	info.Name = utils.Stringp("marvel")
	sketch, err := NewCMLSketch(info)

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	values := [][]byte{
		[]byte("sabertooth"),
		[]byte("thunderbolt"),
		[]byte("havoc"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("havoc")}

	if _, err := sketch.Add(values); err != nil {
		t.Error("expected no errors, got", err)
	}

	if res, err := sketch.Get([][]byte{[]byte("cyclops")}); err != nil {
		t.Error("expected no errors, got", err)
	} else if res.(*pb.FrequencyResult).Frequencies[0].GetCount() != 3 {
		t.Error("expected 'cyclops' count == 3, got", res.(*pb.FrequencyResult).Frequencies[0].GetCount())
	}
}
