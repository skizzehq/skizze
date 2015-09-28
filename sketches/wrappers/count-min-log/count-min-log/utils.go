package cml

import (
	"github.com/dgryski/go-farm"
	"github.com/lazybeaver/xorshift"
)

var rnd = xorshift.NewXorShift64Star(42)

func randFloat() float64 {
	return float64(rnd.Next()%10e5) / 10e5
}

func hash(s []byte, i, w uint) uint {
	return uint(farm.Hash64WithSeed(s, uint64(i))) % w
}
