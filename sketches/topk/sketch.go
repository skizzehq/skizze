package topk

import (
	"bytes"
	"encoding/gob"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/sketches/topk/go-topk"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

const defaultRank = 100.0

// Sketch is the toplevel sketch to control the HLL implementation
type Sketch struct {
	*datamodel.Info
	impl *topk.Stream
}

// ResultElement ...
type ResultElement topk.Element

// NewSketch ...
func NewSketch(info *datamodel.Info) (*Sketch, error) {
	if info.Properties.Rank == 0 {
		info.Properties.Rank = defaultRank
	}
	d := Sketch{info, topk.New(int(info.Properties.Rank))}

	return &d, nil
}

// Add ...
func (d *Sketch) Add(values [][]byte) (bool, error) {
	for _, value := range values {
		str := string(value)
		err := d.impl.Insert(str, 1)
		if err != nil {
			logger.Error.Println("an error has occurred while populating sketch: " + err.Error())
			return false, err
		}
	}
	return true, nil
}

// Get ...
func (d *Sketch) Get(interface{}) (interface{}, error) {
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
func (d *Sketch) Marshal() ([]byte, error) {
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
func (d *Sketch) Unmarshal(info *datamodel.Info, data []byte) error {
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
