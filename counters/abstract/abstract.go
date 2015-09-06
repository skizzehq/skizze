package abstract

/*
Counter ...
*/
type Counter interface {
	Add([]byte) (bool, error)
	AddMultiple([][]byte) (bool, error)
	Remove([]byte) (bool, error)
	RemoveMultiple([][]byte) (bool, error)
	GetCount() uint
	Clear() (bool, error)
	GetType() string
	GetID() string
	GetFrequency([][]byte) interface{}
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

/*
Default		=> HLLPP
Purgable	=> CuckooFilter
Frequency	=> Count-min sketch
Expirable	=> Sliding HLL
*/
const (
	HLLPP = "hllpp"
	CML   = "cml"
	TopK  = "topk"
)
