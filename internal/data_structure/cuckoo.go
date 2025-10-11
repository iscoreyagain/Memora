package data_structure

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"

	"github.com/spaolacci/murmur3"
)

const InitialSeed uint32 = 0x9747b28c
const MAX_KICKOUT int = 500

// The ideal buckets'size of Cuckoo filter is should POWER of 2 -> Makes modulo and XOR math
// faster.
type CuckooFilter struct {
	Items           int
	Error           float64   // The desired error rate - defined by user or in default (0.01%)
	buckets         []*Bucket // The slice of buckets - each bucket contains multiple `fingerprints`
	bucketSize      int       // Number of slots for each bucket
	fingerprintSize uint8     // Bits per fingerprint
	bits            uint64    // The total number of bits
	bytes           uint64    // The total number of bytes (bits / 8)
	Capacity        int       // The total number of buckets
}

// Used for bit-packing
type Bucket struct {
	data      []byte // Packed fingerprints
	SlotBits  uint8  // bits per fingerprint
	SlotCount int    // slots per bucket
}

func (b *Bucket) Insert(fp uint64) bool {
	for i := 0; i < b.SlotCount; i++ {
		if b.GetSlot(i) == 0 {
			b.SetSlot(i, fp)
			return true
		}
	}
	return false
}

func (b *Bucket) SwapSlot(slotIdx int, newFp uint64) uint64 {
	oldFp := b.GetSlot(slotIdx)
	b.SetSlot(slotIdx, newFp)

	return oldFp
}

func (b *Bucket) GetSlot(slotIdx int) uint64 {
	// The start of each fingerprint will be (0 -> fingerprintSize - 1), (fingerprintSize -> 2 * fingerprintSize - 1)...
	offset := slotIdx * int(b.SlotBits)

	// Which one among the array of bytes?
	byteIdx := offset / 8

	// Each byte has 8 bits
	// For .i.e we have 5-bit fingerprint's size and array of bytes is 3 bytes
	// So that -> the first fingerprint will use the first 5 bits (index 0 -> 4) of the first byte
	// then the second fingerprint will use the last 3 bits of the first byte + 2 bits of the 2nd
	// The third would go from the index 2 -> 6 of the 2nd
	// Lastly, the last slot will hold the last bit of the 2nd + 4 first bits of the 3rd
	bitIdx := offset % 8

	bitsLeft := b.SlotBits
	var res uint64
	var shift uint

	for bitsLeft > 0 {
		curByte := b.data[byteIdx]
		remainingBits := 8 - bitIdx
		bitsToRead := bitsLeft

		// If the bits in that byte pos is smaller than our bits need to write/flip:
		// [7][6][5][4][3][2][1] (bits) -> [0][0][0][1][0][1][1][0]
		if int(bitsToRead) > remainingBits {
			bitsToRead = uint8(remainingBits)
		}

		mask := uint64((1 << bitsToRead) - 1)
		extracted := uint64(curByte>>bitIdx) & mask

		res |= extracted << uint64(shift)

		bitsLeft -= bitsToRead
		shift += uint(bitsToRead)
		byteIdx++
		bitIdx = 0
	}

	return res
}

func (b *Bucket) SetSlot(slotIdx int, fp uint64) {
	offset := slotIdx * int(b.SlotBits)
	byteIdx := offset / 8
	bitIdx := offset % 8

	bitsLeft := b.SlotBits

	for bitsLeft > 0 {
		remainingBits := 8 - bitIdx
		bitsToFlip := bitsLeft

		if int(bitsToFlip) > remainingBits {
			bitsToFlip = uint8(remainingBits)
		}

		// mask to clear all the bits that we are going to overwrite
		mask := uint64((1 << bitsToFlip) - 1)

		// clear the bits
		b.data[byteIdx] &^= byte(mask << bitIdx)

		// start to flip the corresponding bits
		b.data[byteIdx] |= byte((fp << bitIdx) & mask)

		fp >>= bitsToFlip
		bitsLeft -= bitsToFlip
		byteIdx++
		bitIdx = 0
	}
}

func calcFingerSz(errRate float64, bucketSize int) uint8 {
	return uint8(math.Ceil(math.Log2(2 * float64(bucketSize) / errRate)))
}

func calcHash(key string) uint64 {
	hasher := murmur3.New128WithSeed(InitialSeed)
	hasher.Write([]byte(key))
	sum := hasher.Sum(nil)

	hval := binary.LittleEndian.Uint64(sum[:8])

	return hval
}

// The fingerprint for the Cuckoo will be derived from first 8-byte hash value % (1 << size)
// Or: hval & ((1 << size) - 1)
func calcFingerprint(hval uint64, size uint8) uint64 {
	if size == 0 {
		return 0
	}

	if size >= 64 {
		if hval == 0 {
			return 1
		}
		return hval
	}
	mask := uint64((1 << size) - 1)
	fp := hval & mask
	if fp == 0 {
		fp = 1
	}
	return fp
}

func CreateCuckooFilter(errRate float64, capacity int) *CuckooFilter {
	// Calculate the fingerprint's size
	fpSz := calcFingerSz(errRate, 4)
	bits := uint64(capacity * 4 * int(fpSz))
	bytes := (bits + 7) / 8

	buckets := make([]*Bucket, capacity)
	ideal := uint8(math.Ceil(float64((4 * fpSz)) / 8))
	for i := 0; i < capacity; i++ {
		buckets[i] = &Bucket{
			data:      make([]byte, ideal),
			SlotBits:  fpSz,
			SlotCount: 4,
		}
	}

	return &CuckooFilter{
		Error:           errRate,
		buckets:         buckets,
		bucketSize:      4,
		fingerprintSize: fpSz,
		bits:            bits,
		bytes:           bytes,
		Capacity:        capacity,
	}
}

func (cf *CuckooFilter) Add(entry string) bool {
	hval := calcHash(entry)
	fp := calcFingerprint(hval, cf.fingerprintSize)
	t1 := hval % uint64(cf.Capacity)
	t2 := (t1 ^ (calcHash(fmt.Sprintf("%d", fp)) % uint64(cf.Capacity))) % uint64(cf.Capacity)

	// try to insert into idx `t1` first
	if cf.buckets[t1].Insert(fp) {
		return true
	}
	// if `t1` has already been used, switch to `t2`
	if cf.buckets[t2].Insert(fp) {
		return true
	}

	i := t1
	for j := 0; j < MAX_KICKOUT; j++ {
		randSlot := rand.Intn(cf.bucketSize)
		oldFp := cf.buckets[i].SwapSlot(randSlot, fp)

		fp = oldFp

		i = (i ^ (calcHash(fmt.Sprintf("%d", fp)) % uint64(cf.Capacity))) % uint64(cf.Capacity)
		if cf.buckets[i].Insert(fp) {
			return true
		}
	}
	return false
}

func (cf *CuckooFilter) Exist(entry string) bool {
	hval := calcHash(entry)
	fp := calcFingerprint(hval, cf.fingerprintSize)
	t1 := hval % uint64(cf.Capacity)
	t2 := (t1 ^ (calcHash(fmt.Sprintf("%d", fp)) % uint64(cf.Capacity))) % uint64(cf.Capacity)

	for _, idx := range []uint64{t1, t2} {
		bucket := cf.buckets[idx]
		for i := 0; i < bucket.SlotCount; i++ {
			if bucket.GetSlot(i) == fp {
				return true
			}
		}
	}
	return false
}

func (cf *CuckooFilter) Remove(entry string) bool {
	hval := calcHash(entry)
	fp := calcFingerprint(hval, cf.fingerprintSize)
	t1 := hval % uint64(cf.Capacity)
	t2 := (t1 ^ (calcHash(fmt.Sprintf("%d", fp)) % uint64(cf.Capacity))) % uint64(cf.Capacity)

	for _, idx := range []uint64{t1, t2} {

		bucket := cf.buckets[idx]
		for i := 0; i < bucket.SlotCount; i++ {
			if bucket.GetSlot(i) == fp {
				bucket.SetSlot(i, 0)
				return true
			}
		}
	}
	return false
}

func (b *Bucket) IsFull() bool {
	for i := 0; i < b.SlotCount; i++ {
		if b.GetSlot(i) == 0 {
			return false
		}
	}
	return true
}
