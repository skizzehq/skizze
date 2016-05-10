package sketches

import (
	"fmt"
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"math/rand"
	"testutils"
	"utils"
)

func TestAddCML(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1000000)
	info.Name = utils.Stringp("marvel")
	typ := pb.SketchType_FREQ
	info.Type = &typ
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

func TestAddCMLThreshold(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	typ := pb.SketchType_FREQ
	info.Type = &typ
	sketch, err := NewCMLSketch(info)

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	rValues := make(map[string]uint64)
	var sValues [][]byte
	thresholdSize := int64(sketch.threshold.size)
	for i := int64(0); i < info.GetProperties().GetMaxUniqueItems()/10; i++ {
		value := fmt.Sprintf("value-%d", i)
		freq := uint64(rand.Int63()) % 100

		values := make([][]byte, freq, freq)
		for i := range values {
			values[i] = []byte(value)
		}
		if _, err := sketch.Add(values); err != nil {
			t.Error("expected no errors, got", err)
		}
		rValues[value] = freq
		sValues = append(sValues, []byte(value))
		// Threshold should be nil once more than 10% is filled
		if sketch.threshold != nil && i >= thresholdSize {
			t.Error("expected threshold == nil for i ==", i)
		}

		if res, err := sketch.Get(sValues); err != nil {
			t.Error("expected no errors, got", err)
		} else {
			tmp := res.(*pb.FrequencyResult)
			mres := tmp.GetFrequencies()
			for i := 0; i < len(mres); i++ {
				if key := mres[i].GetValue(); rValues[key] != uint64(mres[i].GetCount()) {
					t.Fatalf("expected %s: %d, got %d", mres[i].GetValue(), rValues[key], uint64(mres[i].GetCount()))
				}
			}
		}
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
