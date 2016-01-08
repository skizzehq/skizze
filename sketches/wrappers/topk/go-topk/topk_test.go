package topk

import (
	"bufio"
	"log"
	"os"
	"sort"
	"testing"
)

type freqs struct {
	keys   []string
	counts map[string]int
}

func (f freqs) Len() int { return len(f.keys) }

// Actually 'Greater', since we want decreasing
func (f *freqs) Less(i, j int) bool {
	return f.counts[f.keys[i]] >= f.counts[f.keys[j]] || f.counts[f.keys[i]] == f.counts[f.keys[j]] && f.keys[i] < f.keys[j]
}

func (f *freqs) Swap(i, j int) { f.keys[i], f.keys[j] = f.keys[j], f.keys[i] }

func TestTopK(t *testing.T) {

	f, err := os.Open("testdata/domains.txt")

	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	tk := New(100)
	exact := make(map[string]int)

	for scanner.Scan() {

		item := scanner.Text()

		exact[item]++
		tk.Insert(item, 1)
	}

	if err := scanner.Err(); err != nil {
		log.Println("error during scan: ", err)
	}

	var keys []string

	for k := range exact {
		keys = append(keys, k)
	}

	freq := &freqs{keys: keys, counts: exact}

	sort.Sort(freq)

	top := tk.Keys()

	// at least the top 25 must be in order
	for i := 0; i < 25; i++ {
		if top[i].Key != freq.keys[i] {
			t.Errorf("key mismatch: idx=%d top=%s (%d) exact=%s (%d)", i, top[i].Key, top[i].Count, freq.keys[i], freq.counts[freq.keys[i]])
		}
	}
}
