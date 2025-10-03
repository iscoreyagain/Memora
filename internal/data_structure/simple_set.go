package data_structure

// Simple Set is an unordered collections of unique strings (members)

type SimpleSet struct {
	key  string
	dict map[string]struct{}
}

func NewSimpleSet(key string) *SimpleSet {
	return &SimpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

// Implement SADD - Add unique member into the set
// ** Return the number of members that successfully added to the set

func (s *SimpleSet) Add(members ...string) int {
	added := 0
	for _, member := range members {
		if _, exist := s.dict[member]; !exist {
			s.dict[member] = struct{}{}
			added++
		}
	}

	return added
}

// Implement SREM - Remove a/some specified elements from the set
// ** Return the number of elems that were successfully removed from
func (s *SimpleSet) Rem(members ...string) int {
	removed := 0
	for _, member := range members {
		if _, exist := s.dict[member]; exist {
			delete(s.dict, member)
			removed++
		}
	}

	return removed
}

// SISMEMBER
func (s *SimpleSet) IsMember(member string) int {
	if _, exist := s.dict[member]; exist {
		return 1
	}
	return 0
}

// SMEMBERS
func (s *SimpleSet) Members() []string {
	res := make([]string, 0, len(s.dict))

	for key, _ := range s.dict {
		res = append(res, key)
	}

	return res
}

//SCARD
func (s *SimpleSet) Card() int {
	return len(s.dict)
}

// SDIFF
// **Return the set contains the successive key's elements that are not in another sets
func (s *SimpleSet) Difference(sets ...*SimpleSet) *SimpleSet {
	diff := make(map[string]struct{})

	for elem := range s.dict {
		found := false
		for i := 0; i < len(sets); i++ {
			if _, exist := sets[i].dict[elem]; exist {
				found = true
				break
			}
		}
		if !found {
			diff[elem] = struct{}{}
		}
	}

	return &SimpleSet{
		key:  "",
		dict: diff,
	}
}

// SDIFFSTORE
// **Similar to `SDIFF` but it also store that set with a new key = destKey
func (s *SimpleSet) DifferenceStore(destKey string, sets ...*SimpleSet) *SimpleSet {
	storeSet := s.Difference(sets...)
	storeSet.key = destKey

	return storeSet
}

//SINTER
func (s *SimpleSet) Intersection(sets ...*SimpleSet) *SimpleSet {
	sets = append(sets, s)

	base := s.dict
	baseIdx := -1

	for i := range sets {
		if baseIdx == -1 || len(base) > len(sets[i].dict) {
			base = sets[i].dict
			baseIdx = i
		}
	}

	mems := make(map[string]struct{})
	for elem := range base {
		found := true
		for i := 0; i < len(sets); i++ {
			if i == baseIdx {
				continue
			}
			if _, exist := sets[i].dict[elem]; !exist {
				found = false
				break
			}
		}
		if found {
			mems[elem] = struct{}{}
		}
	}

	return &SimpleSet{
		key:  "",
		dict: mems,
	}
}

func (s *SimpleSet) IntersectionStore(destKey string, sets ...*SimpleSet) *SimpleSet {
	interSet := s.Intersection(sets...)
	interSet.key = destKey

	return interSet
}

// SUNION
// **Return
func (s *SimpleSet) Union(sets ...*SimpleSet) *SimpleSet {
	sets = append(sets, s)
	union := make(map[string]struct{})

	for i := 0; i < len(sets); i++ {
		set := sets[i]
		for elem := range set.dict {
			union[elem] = struct{}{}
		}
	}

	return &SimpleSet{
		key:  "",
		dict: union,
	}
}

//SUNIONSTORE
func (s *SimpleSet) UnionStore(destKey string, sets ...*SimpleSet) *SimpleSet {
	unionSet := s.Union(sets...)
	unionSet.key = destKey

	return unionSet
}
