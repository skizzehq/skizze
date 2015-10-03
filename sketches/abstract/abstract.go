package abstract

/*
HLLPP	=> HLLPP
CML		=> Count-min-log sketch
TopK	=> Top-K
*/
const (
	HLLPP = "hllpp"
	CML   = "cml"
	TopK  = "topk"
	Dict  = "dict"
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
