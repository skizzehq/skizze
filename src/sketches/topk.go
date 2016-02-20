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
	d := TopKSketch{info, topk.New(int(info.Properties.GetSize()))}
	return &d, nil
}

// Add ...
func (d *TopKSketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		str := string(value)
		d.impl.Insert(str, 1)

	}
	return true, nil
}

// Get ...
func (d *TopKSketch) Get(interface{}) (interface{}, error) {
	keys := d.impl.Keys()
	result := &pb.RankingsResult{
		Rankings: make([]*pb.Rank, len(keys), len(keys)),
	}
	for i, k := range keys {
		result.Rankings[i] = &pb.Rank{
			Value: utils.Stringp(k.Key),
			Count: utils.Int64p(int64(k.Count)),
		}
	}
	return result, nil
}
