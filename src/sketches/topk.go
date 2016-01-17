package sketches

import (
	"bytes"
	"encoding/gob"

	"github.com/dgryski/go-topk"

	"datamodel"
)

const defaultRank = 100.0

// TopKSketch is the toplevel sketch to control the HLL implementation
type TopKSketch struct {
	*datamodel.Info
	impl *topk.Stream
}

// ResultElement ...
type ResultElement topk.Element

// NewTopKSketch ...
func NewTopKSketch(info *datamodel.Info) (*TopKSketch, error) {
	if info.Properties.Rank == 0 {
		info.Properties.Rank = defaultRank
	}
	d := TopKSketch{info, topk.New(int(info.Properties.Rank))}

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
	result := make([]*datamodel.Element, len(keys), len(keys))
	for i, k := range keys {
		result[i] = &datamodel.Element{
			Key:   k.Key,
			Count: k.Count,
			Error: k.Error,
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
		logger.Error.Println("an error has occurred while loading sketch from data: " + err.Error())
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
