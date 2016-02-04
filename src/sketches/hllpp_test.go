package sketches

import (
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

func TestAddHLLPP(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

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

func TestStressHLLPP(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}

	for i := 0; i < 1024; i++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.MaxUniqueItems = utils.Int64p(1024)
		info.Name = utils.Stringp("marvel" + strconv.Itoa(i))

		sketch, err := NewHLLPPSketch(info)

		if err != nil {
			t.Error("expected avengers to have no error, got", err)
		}

		if _, err := sketch.Add(values); err != nil {
			t.Error("expected no errors, got", err)
		}
	}
}
