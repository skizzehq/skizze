package abstract

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

/*
Default		=> HLLPP
Purgable	=> CuckooFilter
Frequency	=> Count-min sketch
Expirable	=> Sliding HLL
RealCount   => Simple real map counter using channels
*/
const (
	HLLPP     = "hllpp"
	CML       = "cml"
	TopK      = "topk"
	RealCount = "realcount"
)
