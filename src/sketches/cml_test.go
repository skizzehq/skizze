package sketches

import (
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
	"testutils"
)

func TestAdd(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

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

func BenchmarkCML(b *testing.B) {
	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}

	for n := 0; n < b.N; n++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.MaxUniqueItems = utils.Int64p(1000)
		info.Name = utils.Stringp("marvel2")
		sketch, err := NewCMLSketch(info)
		if err != nil {
			b.Error("expected no errors, got", err)
		}
		for i := 0; i < 1000; i++ {
			if _, err := sketch.Add(values); err != nil {
				b.Error("expected no errors, got", err)
			}
		}
	}
}
