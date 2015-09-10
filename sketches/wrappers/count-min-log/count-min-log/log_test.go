package cml

import (
	"strconv"
	"testing"
)

// Ensures that Add adds to the set and Count returns the correct
// approximation.
func TestLog8AddAndCount(t *testing.T) {
	log, _ := NewDefaultSketch8()

	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))
	if count := log.Frequency([]byte("a")); uint(count) != 3 {
		t.Errorf("expected 3, got %d", uint(count))
	}

	if count := log.Frequency([]byte("b")); uint(count) != 2 {
		t.Errorf("expected 2, got %d", uint(count))
	}

	if count := log.Frequency([]byte("c")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.Frequency([]byte("d")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.Frequency([]byte("x")); uint(count) != 0 {
		t.Errorf("expected 0, got %d", uint(count))
	}
}

// Ensures that Add adds to the set and Count returns the correct
// approximation.
func TestLog16AddAndCount(t *testing.T) {
	log, _ := NewDefaultSketch16()

	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))
	if count := log.Frequency([]byte("a")); uint(count) != 3 {
		t.Errorf("expected 3, got %d", uint(count))
	}

	if count := log.Frequency([]byte("b")); uint(count) != 2 {
		t.Errorf("expected 2, got %d", uint(count))
	}

	if count := log.Frequency([]byte("c")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.Frequency([]byte("d")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.Frequency([]byte("x")); uint(count) != 0 {
		t.Errorf("expected 0, got %d", uint(count))
	}
}

// Ensures that Reset restores the sketch to its original state.
func TestLog8Reset(t *testing.T) {
	log, _ := NewDefaultSketch8()
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))

	log.Reset()

	for i := uint(0); i < log.k; i++ {
		for j := uint(0); j < log.w; j++ {
			if x := log.store[i][j]; x != 0 {
				t.Errorf("expected matrix to be completely empty, got %d", x)
			}
		}
	}
}

// Ensures that Reset restores the sketch to its original state.
func TestLog16Reset(t *testing.T) {
	log, _ := NewDefaultSketch16()
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))

	log.Reset()

	for i := uint(0); i < log.k; i++ {
		for j := uint(0); j < log.w; j++ {
			if x := log.store[i][j]; x != 0 {
				t.Errorf("expected matrix to be completely empty, got %d", x)
			}
		}
	}
}

func BenchmarkLogAdd(b *testing.B) {
	b.StopTimer()
	log, _ := NewDefaultSketch8()
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		log.IncreaseCount(data[n])
	}
}

func BenchmarkLogCount(b *testing.B) {
	b.StopTimer()
	log, _ := NewDefaultSketch8()
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
		log.IncreaseCount([]byte(strconv.Itoa(i)))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		log.Frequency(data[n])
	}
}

// Ensures that Add adds to the set and Count returns the correct
// approximation.
func TestLog16Marshall(t *testing.T) {
	log, _ := NewDefaultSketch16()

	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))

	data, err := log.Marshall()

	if err != nil {
		t.Error("expected no error marshalling, got", err)
	}

	alog, err := Unmarshall16(data)

	if err != nil {
		t.Error("expected no error unmarshalling, got", err)
	}

	if count := alog.Frequency([]byte("a")); uint(count) != 3 {
		t.Errorf("expected 3, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("b")); uint(count) != 2 {
		t.Errorf("expected 2, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("c")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("d")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("x")); uint(count) != 0 {
		t.Errorf("expected 0, got %d", uint(count))
	}
}

// Ensures that Add adds to the set and Count returns the correct
// approximation.
func TestLog8Marshall(t *testing.T) {
	log, _ := NewDefaultSketch8()

	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))

	data, err := log.Marshall()

	if err != nil {
		t.Error("expected no error marshalling, got", err)
	}

	alog, err := Unmarshall8(data)

	if err != nil {
		t.Error("expected no error unmarshalling, got", err)
	}

	if count := alog.Frequency([]byte("a")); uint(count) != 3 {
		t.Errorf("expected 3, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("b")); uint(count) != 2 {
		t.Errorf("expected 2, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("c")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("d")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := alog.Frequency([]byte("x")); uint(count) != 0 {
		t.Errorf("expected 0, got %d", uint(count))
	}
}
