package sketches

import (
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
	"testutils"
)

func TestAddBloom(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	sketch, err := NewBloomSketch(info)

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

	check := map[string]bool{
		"sabertooth":  true,
		"thunderbolt": true,
		"havoc":       true,
		"cyclops":     true,
		"wolverine":   false,
		"iceman":      false,
		"rogue":       false,
		"storm":       false}

	if res, err := sketch.Get(values); err != nil {
		t.Error("expected no errors, got", err)
	} else {
		tmp := res.(*pb.MembershipResult)
		mres := tmp.GetMemberships()
		for key := range check {
			for i := 0; i < len(mres); i++ {
				if mres[i].GetValue() == key &&
					mres[i].GetIsMember() != check[key] {
					t.Error("expected member == "+strconv.FormatBool(check[key])+", got", mres[i].GetIsMember())
				}
			}
		}
	}
}

func TestStressBloom(t *testing.T) {
	testutils.SetupTests()
	defer testutils.TearDownTests()

	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}

	info := datamodel.NewEmptyInfo()
	info.Properties.MaxUniqueItems = utils.Int64p(1024)
	info.Name = utils.Stringp("marvel")
	sketch, err := NewBloomSketch(info)

	for i := 0; i < 1024; i++ {
		if err != nil {
			t.Error("expected avengers to have no error, got", err)
		}

		if _, err := sketch.Add(values); err != nil {
			t.Error("expected no errors, got", err)
		}
	}
}

func BenchmarkBloom(b *testing.B) {
	testutils.SetupTests()
	defer testutils.TearDownTests()
	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}
	for n := 0; n < b.N; n++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.MaxUniqueItems = utils.Int64p(1000)
		info.Name = utils.Stringp("marvel")
		sketch, err := NewBloomSketch(info)
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
