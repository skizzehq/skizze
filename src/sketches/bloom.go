package sketches

import (
	bloom "github.com/AndreasBriese/bbloom"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// BloomSketch is the toplevel Sketch to control the count-min-log implementation
type BloomSketch struct {
	*datamodel.Info
	impl      *bloom.Bloom
	threshold *Dict
}

// NewBloomSketch ...
func NewBloomSketch(info *datamodel.Info) (*BloomSketch, error) {
	// FIXME: We are converting from int64 to uint
	threshold := NewDict(info)
	d := BloomSketch{info, nil, threshold}
	return &d, nil
}

// Add ...
func (d *BloomSketch) Add(values [][]byte) (bool, error) {
	success := true
	dict := make(map[string]uint)
	if d.threshold != nil {
		s, err := d.threshold.Add(values)
		success = s
		if err != nil {
			return false, err
		}
		if !d.threshold.IsFull() {
			return true, nil
		}
		values = d.threshold.Keys()
		d.threshold = nil
		if d.impl == nil {
			sketch := bloom.New(float64(d.Info.Properties.GetMaxUniqueItems()), 4.0)
			d.impl = &sketch
		}
	}

	for _, v := range values {
		dict[string(v)]++
	}
	for v := range dict {
		d.impl.Add([]byte(v))
	}
	// Fixme: return what was added and what not
	return success, nil
}

// Get ...
func (d *BloomSketch) Get(data interface{}) (interface{}, error) {
	if d.threshold != nil {
		return d.threshold.Get(data)
	}

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
