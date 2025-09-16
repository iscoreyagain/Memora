package data_structure

import "math"

const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201
const ABigSeed uint32 = 0x9747b28c

type BloomFilter struct {
	Hashes      int
	Entries     uint64
	bf          []bool
	Error       float64
	bitPerEntry float64
	bits        uint64
	bytes       uint64
}

type HashValue struct {
	x uint64
	y uint64
}

func calcBpe(err float64) float64 {
	num := math.Log(err)
	return math.Abs(-(num / Ln2Square))
}

func CreateBloomFilter(entries uint64, errorRate float64) *BloomFilter {
	bloom := BloomFilter{
		Entries: entries,
		Error:   errorRate,
	}

	// bitPerEntry = bits / entries
	bloom.bitPerEntry = calcBpe(errorRate)
	bits
}
