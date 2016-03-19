package sketches

import (
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
	"testutils"
)

func TestAddHLLPP(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	sketch, err := NewHLLPPSketch(info)

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

	const expectedCardinality int64 = 4

	if res, err := sketch.Get(values); err != nil {
		t.Error("expected no errors, got", err)
	} else {
		tmp := res.(*pb.CardinalityResult)
		mres := tmp.GetCardinality()
		if mres != int64(expectedCardinality) {
			t.Error("expected cardinality == "+strconv.FormatInt(expectedCardinality, 10)+", got", mres)
		}
	}
}

func BenchmarkHLLPP(b *testing.B) {
	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}

	for n := 0; n < b.N; n++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.Size = utils.Int64p(1000)
		info.Name = utils.Stringp("marvel3")
		sketch, err := NewHLLPPSketch(info)
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
