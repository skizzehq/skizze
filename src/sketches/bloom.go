package sketches

import (
	"github.com/willf/bloom"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// BloomSketch is the toplevel Sketch to control the count-min-log implementation
type BloomSketch struct {
	*datamodel.Info
	impl *bloom.BloomFilter
}

// NewBloomSketch ...
func NewBloomSketch(info *datamodel.Info) (*BloomSketch, error) {
	// FIXME: We are converting from int64 to uint
	sketch := bloom.New(uint(info.Properties.GetMaxUniqueItems()), 4)
	d := BloomSketch{info, sketch}
	return &d, nil
}

// Add ...
func (d *BloomSketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		d.impl.Add(value)
	}
	return true, nil
}

// Marshal ...
func (d *BloomSketch) Marshal() ([]byte, error) {
	return d.impl.GobEncode()
}

// Get ...
func (d *BloomSketch) Get(data interface{}) (interface{}, error) {
	values := data.([][]byte)
	tmpRes := make(map[string]*pb.Membership)
	res := &pb.MembershipResult{
		Memberships: make([]*pb.Membership, len(values), len(values)),
	}
	for i, v := range values {
		if r, ok := tmpRes[string(v)]; ok {
			res.Memberships[i] = r
			continue
		}
		res.Memberships[i] = &pb.Membership{
			Value:    utils.Stringp(string(v)),
			IsMember: utils.Boolp(d.impl.Test(v)),
		}
		tmpRes[string(v)] = res.Memberships[i]
	}
	return res, nil
}

// Unmarshal ...
func (d *BloomSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	sketch := &bloom.BloomFilter{}
	if err := sketch.GobDecode(data); err != nil {
		return err
	}
	d.Info = info
	d.impl = sketch
	return nil
}
