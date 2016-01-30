package sketches

import (
	"bytes"
	"encoding/gob"

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

// Marshal ...
func (d *TopKSketch) Marshal() ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	//  Encode (send) the value.
	err := enc.Encode(d.impl)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

// Unmarshal ...
func (d *TopKSketch) Unmarshal(info *datamodel.Info, data []byte) error {
	var network bytes.Buffer //  Stand-in for a network connection
	_, err := network.Write(data)
	if err != nil {
		logger.Errorf("an error has occurred while loading sketch from data: " + err.Error())
		return err
	}
	dec := gob.NewDecoder(&network) //  Will read from network.

	var counter topk.Stream
	err = dec.Decode(&counter)
	if err != nil {
		return err
	}
	d.Info = info
	d.impl = &counter
	return nil
}
