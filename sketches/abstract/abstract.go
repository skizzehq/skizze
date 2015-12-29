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
Properties ...
*/
type Properties struct {
	Capacity uint `json:"capacity"`
}

/*
State ...
*/
type State struct {
	Additions uint `json:"adds"`
	Deletions uint `json:"deletions"`
}

/*
Info ...
*/
type Info struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	State      *State      `json:"state"`
	Properties *Properties `json:"properties"`
}

/*
NewEmptyPropeties ...
*/
func NewEmptyPropeties() *Properties {
	return &Properties{}
}

/*
NewEmptyState ...
*/
func NewEmptyState() *State {
	return &State{
		Additions: 0,
		Deletions: 0,
	}
}

/*
NewEmptyInfo ...
*/
func NewEmptyInfo() *Info {
	return &Info{
		Properties: NewEmptyPropeties(),
		State:      NewEmptyState(),
	}
}
