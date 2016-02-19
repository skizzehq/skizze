package sketches

import (
	bloom "github.com/AndreasBriese/bbloom"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
    "errors"
)

// BloomSketch is the toplevel Sketch to control the count-min-log implementation
type BloomSketch struct {
	*datamodel.Info
	impl *bloom.Bloom
}

// NewBloomSketch ...
func NewBloomSketch(info *datamodel.Info) (*BloomSketch, error) {
	// FIXME: We are converting from int64 to uint
	sketch := bloom.New(float64(info.Properties.GetMaxUniqueItems()), 4.0)
	d := BloomSketch{info, &sketch}
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
	return make([]byte, 0), errors.New("Method not supported")
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
			IsMember: utils.Boolp(d.impl.Has(v)),
		}
		tmpRes[string(v)] = res.Memberships[i]
	}
	return res, nil
}

// Unmarshal ...
func (d *BloomSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	return errors.New("Method not supported")
}
