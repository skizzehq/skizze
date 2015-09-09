// Package topk implements the Filtered Space-Saving TopK streaming algorithm
/*

The original Space-Saving algorithm:
https://icmi.cs.ucsb.edu/research/tech_reports/reports/2005-23.pdf

The Filtered Space-Saving enhancement:
http://www.l2f.inesc-id.pt/~fmmb/wiki/uploads/Work/misnis.ref0a.pdf

This implementation follows the algorithm of the FSS paper, but not the
suggested implementation.  Specifically, we use a heap instead of a sorted list
of monitored items, and since we are also using a map to provide O(1) access on
update also don't need the c_i counters in the hash table.

Licensed under the MIT license.

*/
package topk

import (
	"container/heap"
	"hash/fnv"
	"sort"
)

// Element is a TopK item
type Element struct {
	Key   string
	Count int
	Error int
}

type elementsByCountDescending []Element

func (elts elementsByCountDescending) Len() int { return len(elts) }
func (elts elementsByCountDescending) Less(i, j int) bool {
	return (elts[i].Count >= elts[j].Count) || (elts[i].Count == elts[j].Count && elts[i].Key < elts[i].Key)
}
func (elts elementsByCountDescending) Swap(i, j int) { elts[i], elts[j] = elts[j], elts[i] }

/*
Keys ...
*/
type Keys struct {
	M    map[string]int
	Elts []Element
}

// Implement the container/heap interface

func (tk *Keys) Len() int { return len(tk.Elts) }
func (tk *Keys) Less(i, j int) bool {
	return (tk.Elts[i].Count < tk.Elts[j].Count) || (tk.Elts[i].Count == tk.Elts[j].Count && tk.Elts[i].Error > tk.Elts[j].Error)
}
func (tk *Keys) Swap(i, j int) {

	tk.Elts[i], tk.Elts[j] = tk.Elts[j], tk.Elts[i]

	tk.M[tk.Elts[i].Key] = i
	tk.M[tk.Elts[j].Key] = j
}

/*
Push ...
*/
func (tk *Keys) Push(x interface{}) {
	e := x.(Element)
	tk.M[e.Key] = len(tk.Elts)
	tk.Elts = append(tk.Elts, e)
}

/*
Pop ...
*/
func (tk *Keys) Pop() interface{} {
	var e Element
	e, tk.Elts = tk.Elts[len(tk.Elts)-1], tk.Elts[:len(tk.Elts)-1]

	delete(tk.M, e.Key)

	return e
}

// Stream calculates the TopK elements for a stream
type Stream struct {
	N      int
	K      Keys
	Alphas []int
}

// New returns a Stream estimating the top n most frequent elements
func New(n int) *Stream {
	return &Stream{
		N:      n,
		K:      Keys{M: make(map[string]int), Elts: make([]Element, 0, n)},
		Alphas: make([]int, n*6), // 6 is the multiplicative constant from the paper
	}
}

// Insert adds an element to the stream to be tracked
func (s *Stream) Insert(x string, count int) error {
	h := fnv.New32a()
	_, err := h.Write([]byte(x))
	if err != nil {
		return err
	}
	xhash := int(h.Sum32()) % len(s.Alphas)

	// are we tracking this element?
	if idx, ok := s.K.M[x]; ok {
		s.K.Elts[idx].Count += count
		heap.Fix(&s.K, idx)
		return nil
	}

	// can we track more elements?
	if len(s.K.Elts) < s.N {
		// there is free space
		heap.Push(&s.K, Element{Key: x, Count: count})
		return nil
	}

	if s.Alphas[xhash]+count < s.K.Elts[0].Count {
		s.Alphas[xhash] += count
		return nil
	}

	// replace the current minimum element
	minKey := s.K.Elts[0].Key

	h.Reset()
	_, err = h.Write([]byte(minKey))
	if err != nil {
		return err
	}
	mkhash := int(h.Sum32()) % len(s.Alphas)
	s.Alphas[mkhash] = s.K.Elts[0].Count

	s.K.Elts[0].Key = x
	s.K.Elts[0].Error = s.Alphas[xhash]
	s.K.Elts[0].Count = s.Alphas[xhash] + count

	// we're not longer monitoring minKey
	delete(s.K.M, minKey)
	// but 'x' is as array position 0
	s.K.M[x] = 0

	heap.Fix(&s.K, 0)
	return nil
}

// Keys returns the current estimates for the most frequent elements
func (s *Stream) Keys() []Element {
	elts := append([]Element(nil), s.K.Elts...)
	sort.Sort(elementsByCountDescending(elts))
	return elts
}
