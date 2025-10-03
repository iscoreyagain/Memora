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

/* func (s *SimpleSet) Intersection(sets ...map[string]struct{}) map[string]struct{} {
	if len(sets) == 0 {
		return map[string]struct{}{}
	}

	base := sets[0]
	baseIdx := 0

	for i := 1; i < len(sets); i++ {
		if len(sets[i]) < len(base) {
			base = sets[i]
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

			if _, exist := sets[i][elem]; !exist {
				found = false
				break
			}
		}
		if found {
			mems[elem] = struct{}{}
		}
	}

	return mems
}
*/

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
