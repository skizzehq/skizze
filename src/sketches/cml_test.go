package sketches

import (
	"testing"

	"datamodel"
	"utils"
)

func TestAddCML(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.Capacity = 1000000000
	info.Name = "marvel"
	info.Type = datamodel.HLLPP
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
	} else if res.(map[string]uint)["cyclops"] != 3 {
		t.Error("expected 'cyclops' count == 3, got", res.(map[string]uint)["cyclops"])
	}
}
