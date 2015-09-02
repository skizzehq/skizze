package cml

import (
	"errors"
	"hash/fnv"
	"math"
	"math/rand"
	"strconv"

	"github.com/seiflotfy/skizze/counters/abstract"
)

func hash(s []byte, i uint, w uint) uint {
	hasher := fnv.New32a()
	hasher.Write(s)
	hasher.Write([]byte(strconv.Itoa(int(i))))
	return uint(hasher.Sum32()) % w
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

	store      *registers
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
		depth = uint(math.Ceil(math.Log(1 / delta)))
	)
	return NewSketch16(id, width, depth, true, 1.00026, true, true, 16)
}

/*
NewSketch16 returns a new Count-Min-Log sketch with 16-bit registers
*/
func NewSketch16(id string, w uint, k uint, conservative bool, exp float64,
	maxSample bool, progressive bool, nBits uint) (*Sketch16, error) {
	store := newRegisters(id, k, w)
	//store2 := newRegisters("id", k, w)
	cMax := math.Pow(2.0, float64(nBits)) - 1.0
	if cMax > math.MaxUint16 {
		return nil,
			errors.New("using 16 bit registers allows a max nBits value of 16")
	}
	return &Sketch16{
		w:            w,
		k:            k,
		conservative: conservative,
		exp:          exp,
		maxSample:    maxSample,
		progressive:  progressive,
		nBits:        nBits,
		store:        store,
		totalCount:   0.0,
		cMax:         cMax,
		id:           id,
	}, nil
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
	sketch, err := NewSketch16(info.ID, uint(w), 7, true, 1.00026, true, true, 16)
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
func (sk *Sketch16) Reset() {
	sk.store = newRegisters(sk.id, sk.k, sk.w)
	sk.store.reset()
	sk.totalCount = 0
}

/*
IncreaseCount increases the count of `s` by one, return true if added and the current count of `s`
*/
func (sk *Sketch16) IncreaseCount(s []byte) (bool, float64) {
	sk.totalCount++
	v := make([]uint16, sk.k)
	vmin := uint16(math.MaxUint16)
	vmax := uint16(0)
	for i := range v {
		//v[i] = sk.store[i][hash(s, uint(i), sk.w)]
		v[i], _ = sk.store.get(uint(i), hash(s, uint(i), sk.w))
		if v[i] < vmin {
			vmin = v[i]
		}
		if v[i] > vmax {
			vmax = v[i]
		}
	}

	var c uint16
	if sk.maxSample {
		c = vmax
	} else {
		c = vmin
	}

	if float64(c) > sk.cMax {
		return false, 0.0
	}

	increase := sk.randomLog(c, 0.0)
	if increase {
		for i := uint(0); i < sk.k; i++ {
			nc := v[i]
			if !sk.conservative || vmin == nc {
				sk.store.set(uint(i), hash(s, i, sk.w), nc+1)
			}
		}
		return increase, fullValue16(vmin+1, sk.getExp(vmin+1))
	}
	return false, fullValue16(vmin, sk.getExp(vmin))
}

/*
GetCount returns the count of `s`
*/
func (sk *Sketch16) GetCount(s []byte) float64 {
	clmin := uint16(math.MaxUint16)
	for i := uint(0); i < sk.k; i++ {
		cl, _ := sk.store.get(uint(i), hash(s, uint(i), sk.w))
		if cl < clmin {
			clmin = cl
		}
	}
	c := clmin
	return fullValue16(c, sk.getExp(c))
}

/*
GetProbability returns the error probability of `s`
*/
func (sk *Sketch16) GetProbability(s []byte) float64 {
	v := sk.GetCount(s)
	if v > 0 {
		return v / float64(sk.totalCount)
	}
	return 0
}

/*
Count returns the total count
*/
func (sk *Sketch16) Count() uint {
	return sk.totalCount
}
