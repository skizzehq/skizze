package datamodel

import pb "datamodel/protobuf"

/*
HLLPP	=> HyperLogLogPlusPlus
CML		=> Count-min-log sketch
TopK	=> Top-K
Bloom 	=> Bloom Filter
*/
const (
	DOM   = "dom"
	HLLPP = "card"
	CML   = "freq"
	TopK  = "rank"
	Bloom = "memb"
)

/*
  MEMB = 1;
  FREQ = 2;
  RANK = 3;
  CARD = 4;
*/
var typeMap = map[pb.SketchType]string{
	pb.SketchType_MEMB: Bloom,
	pb.SketchType_FREQ: CML,
	pb.SketchType_RANK: TopK,
	pb.SketchType_CARD: HLLPP,
}

// GetTypes ...
func GetTypes() []string {
	return []string{HLLPP, CML, TopK, Bloom}
}

// GetTypeString ...
func GetTypeString(typ pb.SketchType) string {
	return typeMap[typ]
}

// GetTypesPb ...
func GetTypesPb() []pb.SketchType {
	return []pb.SketchType{
		pb.SketchType_MEMB,
		pb.SketchType_FREQ,
		pb.SketchType_RANK,
		pb.SketchType_CARD,
	}
}
