package cuckoofilter

import "github.com/seiflotfy/skizze/storage"

const fingerprintSize = 1
const bucketSize = 4

type fingerprint [fingerprintSize]byte
type bucket [bucketSize]fingerprint

var nullFp = fingerprint{0}

func (b *bucket) insert(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == nullFp {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b *bucket) delete(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == fp {
			b[i] = nullFp
			return true
		}
	}
	return false
}

func (b *bucket) getFingerprintIndex(fp fingerprint) int {
	for i, tfp := range b {
		if tfp == fp {
			return i
		}
	}
	return -1
}

type buckets struct {
	id         string
	numBuckets uint
}

func newBuckets(id string, numBuckets uint) (*buckets, error) {
	storageManager := storage.GetManager()
	storageManager.Create(id)
	bs := &buckets{id, numBuckets}
	for i := uint(0); i < numBuckets; i++ {
		err := bs.setBucket(i, &bucket{})
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}

func (bs *buckets) getBucket(i uint) (*bucket, error) {
	storageManager := storage.GetManager()
	bytes, err := storageManager.LoadData(bs.id, int64(i*bucketSize*fingerprintSize), int64(bucketSize*fingerprintSize))
	if err != nil {
		return nil, err
	}
	bckt := &bucket{}
	for i, b := range bytes {
		bckt[i/fingerprintSize][i%fingerprintSize] = b
	}
	return bckt, nil
}

func (bs *buckets) setBucket(i uint, bckt *bucket) error {
	storageManager := storage.GetManager()
	data := make([]byte, bucketSize*fingerprintSize, bucketSize*fingerprintSize)
	for i, fp := range bckt {
		for j, v := range fp {
			data[i*fingerprintSize+j] = v
		}
	}
	storageManager.SaveData(bs.id, data, int64(i*bucketSize*fingerprintSize))
	return nil
}
