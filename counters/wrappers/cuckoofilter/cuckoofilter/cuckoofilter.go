package cuckoofilter

import "math/rand"

const maxCuckooCount = 500

/*
CuckooFilter represents a probabalistic counter
*/
type CuckooFilter struct {
	ID    string
	bs    *buckets
	count uint
}

/*
NewCuckooFilter returns a new cuckoofilter with a given capacity
*/
func NewCuckooFilter(ID string, capacity uint) *CuckooFilter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	//FIXME: return error
	bs, _ := newBuckets(ID, capacity)
	return &CuckooFilter{ID, bs, 0}
}

/*
NewDefaultCuckooFilter returns a new cuckoofilter with the default capacity of 1000000
*/
func NewDefaultCuckooFilter(ID string) *CuckooFilter {
	return NewCuckooFilter(ID, 1000000)
}

/*
Lookup returns true if data is in the counter
*/
func (cf *CuckooFilter) Lookup(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bs.numBuckets)
	// FIXME: deal with errors
	b1, _ := cf.bs.getBucket(i1)
	b2, _ := cf.bs.getBucket(i2)
	return b1.getFingerprintIndex(fp) > -1 || b2.getFingerprintIndex(fp) > -1
}

/*
Insert inserts data into the counter and returns true upon success
*/
func (cf *CuckooFilter) Insert(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bs.numBuckets)
	if cf.insert(fp, i1) || cf.insert(fp, i2) {
		return true
	}
	return cf.reinsert(fp, i2)
}

/*
InsertUnique inserts data into the counter if not exists and returns true upon success
*/
func (cf *CuckooFilter) InsertUnique(data []byte) bool {
	if cf.Lookup(data) {
		return false
	}
	return cf.Insert(data)
}

func (cf *CuckooFilter) insert(fp fingerprint, i uint) bool {
	// FIXME: return error
	b, _ := cf.bs.getBucket(i)
	if b.insert(fp) {
		cf.bs.setBucket(i, b)
		cf.count++
		return true
	}
	return false
}

func (cf *CuckooFilter) reinsert(fp fingerprint, i uint) bool {
	for k := 0; k < maxCuckooCount; k++ {
		j := rand.Intn(bucketSize)
		oldfp := fp

		// FIXME: return error
		b, _ := cf.bs.getBucket(i)

		fp = b[j]
		b[j] = oldfp
		cf.bs.setBucket(i, b)

		// look in the alternate location for that random element
		i = getAltIndex(fp, i, cf.bs.numBuckets)
		if cf.insert(fp, i) {
			return true
		}
	}
	return false
}

/*
Delete data from counter if exists and return if deleted or not
*/
func (cf *CuckooFilter) Delete(data []byte) bool {
	i1, i2, fp := getIndicesAndFingerprint(data, cf.bs.numBuckets)
	return cf.delete(fp, i1) || cf.delete(fp, i2)
}

func (cf *CuckooFilter) delete(fp fingerprint, i uint) bool {
	b, _ := cf.bs.getBucket(i)
	if b.delete(fp) {
		cf.bs.setBucket(i, b)
		cf.count--
		return true
	}
	return false
}

/*
GetCount returns the number of items in the counter
*/
func (cf *CuckooFilter) GetCount() uint {
	return cf.count
}
