package sketches

import (
	"datamodel"
	pb "datamodel/protobuf"
	"fmt"
	"utils"
)

// Dict ...
type Dict struct {
	*datamodel.Info
	impl map[string]uint
	size uint
}

// NewDict ...
func NewDict(info *datamodel.Info) *Dict {
	sketch := make(map[string]uint)
	size := info.GetProperties().GetMaxUniqueItems() / 10
	d := Dict{info, sketch, uint(size)}
	return &d
}

// IsFull ...
func (d *Dict) IsFull() bool {
	return uint(len(d.impl)) >= d.size
}

// Add ...
func (d *Dict) Add(values [][]byte) (bool, error) {
	for _, v := range values {
		d.impl[string(v)]++
	}
	return true, nil
}

// Keys ...
func (d *Dict) Keys() [][]byte {
	keys := make([][]byte, len(d.impl), len(d.impl))
	i := 0
	for k := range d.impl {
		keys[i] = []byte(k)
		i++
	}
	return keys
}

// Get ...
func (d *Dict) Get(data interface{}) (interface{}, error) {
	typ := d.Info.Sketch.GetType()
	switch datamodel.GetTypeString(typ) {
	case datamodel.Bloom:
		return d.getMemb(data)
	case datamodel.CML:
		return d.getFreq(data)
	}
	return nil, fmt.Errorf("Unknown error: %v", d.Info.GetType().String()) // FIXME: return some error
}

func (d *Dict) getMemb(data interface{}) (interface{}, error) {
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
		_, ok := d.impl[string(v)]
		res.Memberships[i] = &pb.Membership{
			Value:    utils.Stringp(string(v)),
			IsMember: utils.Boolp(ok),
		}
		tmpRes[string(v)] = res.Memberships[i]
	}
	return res, nil
}

func (d *Dict) getFreq(data interface{}) (interface{}, error) {
	fmt.Println("----->")
	values := data.([][]byte)
	res := &pb.FrequencyResult{
		Frequencies: make([]*pb.Frequency, len(values), len(values)),
	}
	tmpRes := make(map[string]*pb.Frequency)
	for i, v := range values {
		if r, ok := tmpRes[string(v)]; ok {
			res.Frequencies[i] = r
			continue
		}
		res.Frequencies[i] = &pb.Frequency{
			Value: utils.Stringp(string(v)),
			Count: utils.Int64p(int64(d.impl[string(v)])),
		}
		tmpRes[string(v)] = res.Frequencies[i]
	}
	return res, nil
}
