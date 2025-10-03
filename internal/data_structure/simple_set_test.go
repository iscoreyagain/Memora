package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSet_SINTER(t *testing.T) {
	// Setup sets
	s1 := NewSimpleSet("s1")
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")

	s2 := NewSimpleSet("s2")
	s2.Add("b")
	s2.Add("c")
	s2.Add("d")

	s3 := NewSimpleSet("s3")
	s3.Add("c")
	s3.Add("d")
	s3.Add("e")

	// Test 1: intersection of s1, s2, s3 => {c}
	result := s1.Intersection(s2, s3)
	members := result.Members()
	assert.Len(t, members, 1)
	assert.Contains(t, members, "c")

	// Test 2: intersection of s1, s2 => {b, c}
	result2 := s1.Intersection(s2)
	members2 := result2.Members()
	assert.Len(t, members2, 2)
	assert.Contains(t, members2, "b")
	assert.Contains(t, members2, "c")

	// Test 3: intersection with missing set => empty
	emptySet := NewSimpleSet("empty")
	result3 := s1.Intersection(emptySet)
	members3 := result3.Members()
	assert.Len(t, members3, 0)

	// Test 4: intersection of single set => itself
	result4 := s1.Intersection()
	members4 := result4.Members()
	assert.Len(t, members4, 3)
	assert.Contains(t, members4, "a")
	assert.Contains(t, members4, "b")
	assert.Contains(t, members4, "c")
}
