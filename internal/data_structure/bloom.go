package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201
const ABigSeed uint32 = 0x9747b28c

type BloomFilter struct {
	Hashes      int     // The number of hash functions
	Entries     uint64  // The expected ~ maximum entries that BloomFilter can hold
	bf          []uint8 // The actual bit array
	Error       float64 // Target false positive rate
	bitPerEntry float64 // The number of bits allocated for each entry
	bits        uint64  // The total number of bits
	bytes       uint64  // The total number of bytes (bits / 8)
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
	bits := uint64(float64(entries) * bloom.bitPerEntry)

	if bits%64 != 0 {
		bloom.bytes = ((bits / 64) + 1) * 8
	} else {
		bloom.bytes = bits / 8
	}
	bloom.Hashes = int(math.Ceil(Ln2 * bloom.bitPerEntry))
	bloom.bits = bloom.bytes * 8
	bloom.bf = make([]uint8, bloom.bytes)

	return &bloom
}

func (b *BloomFilter) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x, y := hasher.Sum128()

	return HashValue{
		x: x,
		y: y,
	}
}

func (b *BloomFilter) Add(entry string) {
	var hash, bytePos uint64
	initHash := b.CalcHash(entry)
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.x + initHash.y*uint64(i)) % b.bits
		bytePos = hash >> 3
		b.bf[bytePos] |= 1 << (hash % 8)
	}
}

func (b *BloomFilter) Exist(entry string) bool {
	var hash, bytePos uint64
	initHash := b.CalcHash(entry)
	for i := 0; i < b.Hashes; i++ {
		hash = (initHash.x + initHash.y*uint64(i)) & b.bits
		bytePos = hash >> 3
		if (b.bf[bytePos] & (1 << (hash % 8))) == 0 {
			return false
		}
	}
	return true
}
