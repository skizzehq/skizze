package abstract

/*
Counter ...
*/
type Counter interface {
	Add([]byte) bool
	AddMultiple([][]byte) bool
	Remove([]byte) bool
	RemoveMultiple([][]byte) bool
	GetCount() uint
	Clear() bool
}

/*
Info ...
*/
type Info struct {
	ID   string
	Type string
}
