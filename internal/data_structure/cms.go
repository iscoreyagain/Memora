package data_structure

// For further information or to check the integrity of math-proof behind these formula, please check this reference:
// https://en.wikipedia.org/wiki/Count%E2%80%93min_sketch

import (
	"math"

	"github.com/spaolacci/murmur3"
)

type CMS struct {
	rows    uint64
	columns uint64
	counter [][]uint64
}

func CreateCMS(r, c uint64) *CMS {
	cms := &CMS{
		rows:    r,
		columns: c,
	}

	cms.counter = make([][]uint64, r)
	for i := uint64(0); i < r; i++ {
		cms.counter[i] = make([]uint64, c)
	}

	return cms
}

func CalcCMSDim(errRate, errProb float64) (uint64, uint64) {
	rows := uint64(math.Ceil(math.E / errRate))
	cols := uint64(math.Ceil(math.Log(1.0 / errProb)))

	return rows, cols
}

func (c *CMS) calcHash(item string, seed uint32) uint64 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))

	return uint64(hasher.Sum32())
}

func (c *CMS) IncrBy(item string, value uint64) uint64 {
	var minCount uint64 = math.MaxUint64

	for i := 0; i < int(c.rows); i++ {
		hash := c.calcHash(item, uint32(i))

		j := hash % c.columns

		c.counter[i][j] += value

		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}
	return minCount
}

func (c *CMS) Estimate(item string) uint64 {
	var count uint64 = math.MaxUint64

	for i := 0; i < int(c.rows); i++ {
		hash := c.calcHash(item, uint32(i))
		j := hash % c.columns
		if c.counter[i][j] < count {
			count = c.counter[i][j]
		}
	}
	return count
}
