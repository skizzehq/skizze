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
}

/*
Info ...
*/
type Info struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Capacity uint64 `json:"capacity"`
}

/*
Default		=> HLLPP
Purgable	=> CuckooFilter
Frequency	=> Count-min sketch
Expirable	=> Sliding HLL
*/
const (
	Default   = "default"
	Purgable  = "purgable"
	Frequency = "frequency"
	Expiring  = "expiring"
)
