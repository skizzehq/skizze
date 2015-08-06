package counters

/*
Immutable ...
Mutable ...
*/
const (
	Immutable string = "immutable"
	Mutable          = "mutable"
)

/*
AbstractCounter ...
*/
type AbstractCounter struct {
	ID        string
	Type      string
	Capacity  uint
	Prob      float64
	SliceSize uint
	Count     uint
}

/*
GetCount ...
*/
func (cs *AbstractCounter) GetCount() uint {
	return cs.Count
}

/*
Clear ...
*/
func (cs *AbstractCounter) Clear() error {
	return nil
}

/*
Purge ...
*/
func (cs *AbstractCounter) Purge() error {
	return nil
}

/*
Flush ...
*/
func (cs *AbstractCounter) Flush() error {
	return nil
}
