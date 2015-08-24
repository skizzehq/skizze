package abstract

/*
Counter ...
*/
type Counter interface {
	Add([]byte) (bool, error)
	AddMultiple([][]byte) (bool, error)
	Remove([]byte) (bool, error)
	RemoveMultiple([][]byte) (bool, error)
	GetCount() interface{}
	Clear() (bool, error)
	GetType() string
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
	Default             = "default"
	Cardinality         = "cardinality"
	PurgableCardinality = "pcardinality"
	Frequency           = "frequency"
	PurgableFrequency   = "pfrequency"
	Expiring            = "expiring"
	TopK                = "topk"
)
