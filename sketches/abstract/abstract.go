package abstract

/*
HLLPP	=> HyperLogLogPlusPlus
CML		=> Count-min-log sketch
TopK	=> Top-K
Dict  => dictionary
Bloom => Bloom Filter
*/
const (
	HLLPP = "hllpp"
	CML   = "cml"
	TopK  = "topk"
	Dict  = "dict"
	Bloom = "bloom"
)

/*
Sketch ...
*/
type Sketch interface {
	Add([]byte) (bool, error)
	AddMultiple([][]byte) (bool, error)
	Remove([]byte) (bool, error)
	RemoveMultiple([][]byte) (bool, error)
	GetCount() uint
	Clear() (bool, error)
	GetFrequency([][]byte) interface{}
	Marshal() ([]byte, error)
}

/*
Info ...
*/
type Info struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	State      map[string]uint64  `json:"state"`
	Properties map[string]float64 `json:"properties"`
}
