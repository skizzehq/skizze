package cml

import (
	"errors"
	"math"
	"math/rand"

	"github.com/dgryski/go-farm"

	"github.com/seiflotfy/skizze/sketches/abstract"
)

func hash(s []byte, i uint, w uint) (uint, error) {
	v := farm.Hash64WithSeed(s, uint64(i))
	return uint(v) % w, nil
}

func value16(c uint16, exp float64) float64 {
	if c == 0 {
		return 0.0
	}
	return math.Pow(exp, float64(c-1))
}

func fullValue16(c uint16, exp float64) float64 {
	if c <= 1 {
		return value16(c, exp)
	}
	return (1.0 - value16(c+1, exp)) / (1.0 - exp)
}

/*
Sketch16 is a Count-Min-Log sketch 16-bit registers
*/
type Sketch16 struct {
	w            uint
	k            uint
	conservative bool
	exp          float64
	maxSample    bool
	progressive  bool
	nBits        uint

	store      [][]uint16 //*registers
	totalCount uint
	cMax       float64
	id         string
}

/*
NewSketch16ForEpsilonDelta ...
*/
func NewSketch16ForEpsilonDelta(id string, epsilon, delta float64) (*Sketch16, error) {
	var (
		width = uint(math.Ceil(math.E / epsilon))
		depth = uint(math.Ceil(math.Log(1 / (1 - delta))))
	)
	return NewSketch16(id, width, depth, true, 1.00026, true, true, 16)
}

/*
NewSketch16 returns a new Count-Min-Log sketch with 16-bit registers
*/
func NewSketch16(id string, w uint, k uint, conservative bool, exp float64,
	maxSample bool, progressive bool, nBits uint) (*Sketch16, error) {

	cMax := math.Pow(2.0, float64(nBits)) - 1.0
	if cMax > math.MaxUint16 {
		return nil,
			errors.New("using 16 bit registers allows a max nBits value of 16")
	}
	sketch := &Sketch16{
		w:            w,
		k:            k,
		conservative: conservative,
		exp:          exp,
		maxSample:    maxSample,
		progressive:  progressive,
		nBits:        nBits,
		totalCount:   0.0,
		cMax:         cMax,
		id:           id,
		//registers:    newRegisters(id, k, w),
	}
	//store, _ := sketch.registers.load()

	//if len(store) == 0 {
	store := make([][]uint16, k, k)
	for i := range store {
		store[i] = make([]uint16, w, w)
	}
	//}
	sketch.store = store
	return sketch, nil
}

/*
NewForCapacity16 returns a new Count-Min-Log sketch with 16-bit registers optimized for a given max capacity and expected error rate
*/
func NewForCapacity16(info *abstract.Info, capacity uint64, e float64) (*Sketch16, error) {
	// e = 2n/w    ==>    w = 2n/e
	if !(e >= 0.001 && e < 1.0) {
		return nil, errors.New("e needs to be >= 0.001 and < 1.0")
	}
	w := float64(2*capacity) / (e * 7)
	sketch, err := NewSketch16(info.ID, uint(w), 2, true, 1.00026, true, true, 16)
	return sketch, err
}

func (sk *Sketch16) randomLog(c uint16, exp float64) bool {
	pIncrease := 1.0 / (fullValue16(c+1, sk.getExp(c+1)) - fullValue16(c, sk.getExp(c)))
	return rand.Float64() < pIncrease
}

func (sk *Sketch16) getExp(c uint16) float64 {
	if sk.progressive == true {
		return 1.0 + ((sk.exp - 1.0) * (float64(c) - 1.0) / sk.cMax)
	}
	return sk.exp
}

/*
Reset the Sketch to a fresh state (all counters set to 0)
*/
func (sk *Sketch16) Reset() error {
	sk.store = make([][]uint16, sk.k, sk.k)
	for i := range sk.store {
		sk.store[i] = make([]uint16, sk.w, sk.w)
	}
	//err := sk.registers.save(sk.store)
	sk.totalCount = 0
	return nil
}

func (sk *Sketch16) getV(s []byte) ([]uint16, uint16, uint16, error) {
	v := make([]uint16, sk.k)
	vmin := uint16(math.MaxUint16)
	vmax := uint16(0)
	for i := range v {
		h, err := hash(s, uint(i), sk.w)
		if err != nil {
			return nil, 0, 0, err
		}
		v[i] = sk.store[i][h]
		if v[i] < vmin {
			vmin = v[i]
		}
		if v[i] > vmax {
			vmax = v[i]
		}
	}
	return v, vmin, vmax, nil
}

/*
IncreaseCount increases the count of `s` by one, return true if added and the current count of `s`
*/
func (sk *Sketch16) IncreaseCount(s []byte) (bool, float64, error) {
	sk.totalCount++
	v, vmin, vmax, err := sk.getV(s)
	if err != nil {
		return false, -1, err
	}
	var c uint16
	if sk.maxSample {
		c = vmax
	} else {
		c = vmin
	}

	if float64(c) > sk.cMax {
		return false, 0.0, nil
	}

	increase := sk.randomLog(c, 0.0)
	if increase {
		for i := uint(0); i < sk.k; i++ {
			nc := v[i]
			if !sk.conservative || vmin == nc {
				h, err := hash(s, uint(i), sk.w)
				if err != nil {
					return false, -1, err
				}
				sk.store[i][h] = nc + 1
			}
		}
		return increase, fullValue16(vmin+1, sk.getExp(vmin+1)), nil
	}
	return false, fullValue16(vmin, sk.getExp(vmin)), nil
}

/*
GetCount returns the count of `s`
*/
func (sk *Sketch16) GetCount(s []byte) (float64, error) {
	clmin := uint16(math.MaxUint16)
	for i := uint(0); i < sk.k; i++ {
		h, err := hash(s, uint(i), sk.w)
		if err != nil {
			return -1, err
		}
		cl := sk.store[i][h]
		if cl < clmin {
			clmin = cl
		}
	}
	c := clmin
	return fullValue16(c, sk.getExp(c)), nil
}

/*
GetProbability returns the error probability of `s`
*/
func (sk *Sketch16) GetProbability(s []byte) (float64, error) {
	v, err := sk.GetCount(s)
	if err != nil {
		return -1, err
	}
	if v > 0 {
		return v / float64(sk.totalCount), nil
	}
	return 0, nil
}

/*
Count returns the total count
*/
func (sk *Sketch16) Count() uint {
	return sk.totalCount
}
