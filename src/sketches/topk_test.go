package sketches

import (
	"strconv"
	"testing"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

func TestAddTopK(t *testing.T) {
	utils.SetupTests()
	defer utils.TearDownTests()

	info := datamodel.NewEmptyInfo()
	info.Properties.Size = utils.Int64p(3)
	info.Name = utils.Stringp("marvel")
	sketch, err := NewTopKSketch(info)

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	values := [][]byte{
		[]byte("sabertooth"),
		[]byte("thunderbolt"),
		[]byte("cyclops"),
		[]byte("thunderbolt"),
		[]byte("thunderbolt"),
		[]byte("havoc"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("havoc")}

	if _, err := sketch.Add(values); err != nil {
		t.Error("expected no errors, got", err)
	}

	type RankingsStruct struct {
		Value    string
		Position int64
		Count    int64
	}
	expectedRankings := make([]*RankingsStruct, 4, 4)

	expectedRankings[0] = &RankingsStruct{
		Value:    "cyclops",
		Count:    4,
		Position: 1,
	}

	expectedRankings[1] = &RankingsStruct{
		Value:    "thunderbolt",
		Count:    3,
		Position: 2,
	}

	expectedRankings[2] = &RankingsStruct{
		Value:    "havoc",
		Count:    2,
		Position: 3,
	}

	expectedRankings[3] = &RankingsStruct{
		Value:    "sabertooth",
		Count:    1,
		Position: 4,
	}

	if res, err := sketch.Get(values); err != nil {
		t.Error("expected no errors, got", err)
	} else {
		tmp := res.(*pb.RankingsResult)
		rres := tmp.GetRankings()
		for i := 0; i < len(rres); i++ {
			count := rres[i].GetCount()
			value := rres[i].GetValue()
			for j := 0; j < len(expectedRankings); j++ {
				if expectedRankings[j].Value == value && expectedRankings[j].Count != count && expectedRankings[j].Position != int64(i) {
					t.Error("expected ranking == "+strconv.FormatInt(expectedRankings[j].Position, 10)+", got", count)
				}
			}
		}
	}
}

func BenchmarkTopK(b *testing.B) {
	values := make([][]byte, 10)
	for i := 0; i < 1024; i++ {
		avenger := "avenger" + strconv.Itoa(i)
		values = append(values, []byte(avenger))
	}

	for n := 0; n < b.N; n++ {
		info := datamodel.NewEmptyInfo()
		info.Properties.Size = utils.Int64p(1000)
		info.Name = utils.Stringp("marvel3")
		sketch, err := NewTopKSketch(info)
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
