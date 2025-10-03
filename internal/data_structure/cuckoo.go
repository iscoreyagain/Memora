package data_structure

import "math"

// The ideal buckets'size of Cuckoo filter is should POWER of 2 -> Makes modulo and XOR math
// faster.
type CuckooFilter struct {
	Error           float64   // The desired error rate - defined by user or in default (0.01%)
	buckets         [][]uint8 // The slice of buckets - each contains multiple `fingerprints`
	bucketSize      int       // Number of slots for each bucket
	fingerprintSize uint8     // Bits per fingerprint
	bits            uint64    // The total number of bits
	bytes           uint64    // The total number of bytes (bits / 8)
	Capacity        int       // The total number of buckets
}

func calcFingerprintSz(errRate float64) uint8 {
	return uint8(math.Ceil(math.Log2(1 / errRate)))
}

func NewCuckooFilter(errRate float64, capacity int) *CuckooFilter {
	buckets := make([][]uint8, capacity)
	for i := 0; i < capacity; i++ {
		buckets[i] = make([]uint8, bucketSize)
	}

	return &CuckooFilter{
		buckets:    buckets,
		bucketSize: bucketSize,
		capacity:   capacity,
	}
}
