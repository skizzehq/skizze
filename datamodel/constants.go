package datamodel

/*
HLLPP	=> HyperLogLogPlusPlus
CML		=> Count-min-log sketch
TopK	=> Top-K
Bloom => Bloom Filter
*/
const (
	DOM   = "dom"
	HLLPP = "card"
	CML   = "freq"
	TopK  = "rank"
	Bloom = "memb"
)

// GetTypes ...
func GetTypes() []string {
	return []string{HLLPP, CML, TopK, Bloom}
}
