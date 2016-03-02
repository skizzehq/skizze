package sketches

import (
	"github.com/dgryski/go-topk"

	"datamodel"
	pb "datamodel/protobuf"

	"utils"
)

// TopKSketch is the toplevel sketch to control the HLL implementation
type TopKSketch struct {
	*datamodel.Info
	impl *topk.Stream
}

// ResultElement ...
type ResultElement topk.Element

// NewTopKSketch ...
func NewTopKSketch(info *datamodel.Info) (*TopKSketch, error) {
	size := int(info.Properties.GetSize()) * 2 // For higher precision
	d := TopKSketch{info, topk.New(size)}
	return &d, nil
}

// Add ...
func (d *TopKSketch) Add(values [][]byte) (bool, error) {
	dict := make(map[string]int)
	for _, v := range values {
		dict[string(v)]++
	}
	for v, count := range dict {
		d.impl.Insert(v, count)
	}
	return true, nil
}

// Get ...
func (d *TopKSketch) Get(interface{}) (interface{}, error) {
	keys := d.impl.Keys()
	size := len(keys)
	if size > int(d.Info.Properties.GetSize())/2 {
		size = int(d.Info.Properties.GetSize()) / 2
	}
	result := &pb.RankingsResult{
		Rankings: make([]*pb.Rank, size, size),
	}
	for i := range result.Rankings {
		k := keys[i]
		result.Rankings[i] = &pb.Rank{
			Value: utils.Stringp(k.Key),
			Count: utils.Int64p(int64(k.Count)),
		}
	}
	return result, nil
}
