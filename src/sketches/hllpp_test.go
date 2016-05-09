package sketches

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"testutils"
	"utils"
)

func TestAddHLLPP(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	typ := pb.SketchType_CARD
	info.Type = &typ
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

func TestAddHLLPPThreshold(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	typ := pb.SketchType_CARD
	info.Type = &typ
	sketch, err := NewHLLPPSketch(info)

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	rValues := make(map[string]uint64)
	thresholdSize := int64(sketch.threshold.size)
	for i := int64(0); i < info.GetProperties().GetMaxUniqueItems()/10; i++ {
		freq := uint64(rand.Int63()) % 100

		value := fmt.Sprintf("value-%d", i)

		values := make([][]byte, freq+1, freq+1)
		for i := range values {
			values[i] = []byte(value)
		}
		if _, err := sketch.Add(values); err != nil {
			t.Error("expected no errors, got", err)
		}
		rValues[value] = freq
		// Threshold should be nil once more than 10% is filled
		if sketch.threshold != nil && i >= thresholdSize {
			t.Error("expected threshold == nil for i ==", i)
		}

		if res, err := sketch.Get(nil); err != nil {
			t.Error("expected no errors, got", err)
		} else {
			tmp := res.(*pb.CardinalityResult)
			mres := tmp.GetCardinality()
			if int64(len(rValues)) != mres {
				t.Fatalf("expected cardinality %d, got %d", len(rValues), mres)
			}
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
