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
	GetFrequency([][]byte) interface{}
}

/*
Info ...
*/
type Info struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Capacity uint64            `json:"capacity"`
	State    map[string]uint64 `json:"state"`
}

/*
Default		=> HLLPP
Purgable	=> CuckooFilter
Frequency	=> Count-min sketch
Expirable	=> Sliding HLL
*/
const (
	Default     = "default"
	Cardinality = "cardinality"
	Frequency   = "frequency"
	TopK        = "topk"
)
