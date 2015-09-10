package cml

import (
	"math/rand"

	"github.com/dgryski/go-farm"
)

func randFloat() float64 {
	return rand.Float64()
}

func hash(s []byte, i, w uint) uint {
	return uint(farm.Hash64WithSeed(s, uint64(i))) % w
}
