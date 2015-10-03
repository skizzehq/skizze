package dict

import (
	"bytes"
	"encoding/gob"

	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/utils"
)

var logger = utils.GetLogger()

/*
Dict ...
*/
type Dict struct {
	hash map[string]int
}

func makeDict() (dict *Dict) {
	return &Dict{
		hash: make(map[string]int),
	}
}

/*
Reset ...
*/
func (dict *Dict) Reset() {
	dict.hash = make(map[string]int)
}

/*
IncreaseCount ...
*/
func (dict *Dict) IncreaseCount(name string) {
	dict.hash[name]++
}

/*
DecreaseCount ...
*/
func (dict *Dict) DecreaseCount(name string) {
	dict.hash[name]--
}

/*
Marshal ...
*/
func (dict *Dict) Marshal() ([]byte, error) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	// Encode (send) the value.
	err := enc.Encode(dict)
	if err != nil {
		return nil, err
	}
	return network.Bytes(), nil
}

/*
Sketch ...
*/
type Sketch struct {
	*abstract.Info
	impl *Dict
}

/*
NewSketch ...
*/
func NewSketch(info *abstract.Info) (*Sketch, error) {
	var dict = makeDict()
	d := Sketch{info, dict}
	return &d, nil
}

/*
Add ...
*/
func (d *Sketch) Add(value []byte) (bool, error) {
	name := string(value)
	d.impl.IncreaseCount(name)
	return true, nil
}

/*
AddMultiple ...
*/
func (d *Sketch) AddMultiple(values [][]byte) (bool, error) {
	for _, value := range values {
		name := string(value)
		d.impl.IncreaseCount(name)
	}
	return true, nil
}

/*
Remove ...
*/
func (d *Sketch) Remove(value []byte) (bool, error) {
	name := string(value)
	d.impl.DecreaseCount(name)
	return true, nil
}

/*
RemoveMultiple ...
*/
func (d *Sketch) RemoveMultiple(values [][]byte) (bool, error) {
	for _, value := range values {
		name := string(value)
		d.impl.DecreaseCount(name)
	}
	return true, nil
}

/*
GetCount ...
*/
func (d *Sketch) GetCount() uint {
	return uint(len(d.impl.hash))
}

/*
Clear ...
*/
func (d *Sketch) Clear() (bool, error) {
	d.impl.Reset()
	return true, nil
}

/*
GetFrequency which is basically our hash
*/
func (d *Sketch) GetFrequency(values [][]byte) interface{} {
	return d.impl.hash
}

/*
Marshal ...
*/
func (d *Sketch) Marshal() ([]byte, error) {
	return d.impl.Marshal()
}

/*
Unmarshal ...
*/
func Unmarshal(info *abstract.Info, data []byte) (*Sketch, error) {
	var network bytes.Buffer // Stand-in for a network connection
	_, err := network.Write(data)
	if err != nil {
		logger.Error.Println("an error has occurred while loading sketch from data: " + err.Error())
		return nil, err
	}
	dec := gob.NewDecoder(&network) // Will read from network.

	var counter Dict
	err = dec.Decode(&counter)
	if err != nil {
		return nil, err
	}
	return &Sketch{info, &counter}, nil
}
