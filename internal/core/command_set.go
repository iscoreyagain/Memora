package core

import (
	"errors"

	"github.com/iscoreyagain/Memora/internal/data_structure"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0] // TODO: check key is used by other types or not
	set, exist := setStore[key]
	if !exist {
		set = data_structure.NewSimpleSet(key)
		setStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SREM' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		set = data_structure.NewSimpleSet(key)
		setStore[key] = set
	}
	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMEMBERS' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	return Encode(set.Members(), false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SISMEMBER' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	return Encode(set.IsMember(args[1]), false)
}

func cmdSCARD(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SCARD' command"), false)
	}

	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false) //Not existed
	}
	return Encode(set.Card(), false)
}

func cmdSINTER(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SINTER' command"), false)
	}

	var sets []*data_structure.SimpleSet

	for _, key := range args {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].Intersection(sets[1:]...)

	return Encode(resultSet.Members(), false)
}

func cmdSINTERSTORE(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SINTERSTORE' command"), false)
	}

	var sets []*data_structure.SimpleSet
	destKey := args[0]

	for _, key := range args[1:] {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].IntersectionStore(destKey, sets...)

	// Replace the old key with the new one
	delete(setStore, destKey)

	// Assign the new one
	setStore[destKey] = resultSet

	numbers := resultSet.Card()

	// Return the number of elems in the stored set
	return Encode(numbers, false)
}

func cmdSDIFF(args []string) []byte { // SDIFF key [key...] -> parsing -> ["key", "key1", "key2",...]
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SDIFF' command"), false)
	}

	var sets []*data_structure.SimpleSet

	for _, key := range args {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].Difference(sets[1:]...)

	return Encode(resultSet.Members(), false)
}

func cmdSDIFFSTORE(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SDIFFSTORE' command"), false)
	}

	destKey := args[0]
	var sets []*data_structure.SimpleSet

	for _, key := range args[1:] {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].DifferenceStore(destKey, sets[1:]...)

	// Replace the old key with the new one
	delete(setStore, destKey)

	// Assign the new one
	setStore[destKey] = resultSet

	numbers := resultSet.Card()

	// Return the number of elems in the stored set
	return Encode(numbers, false)
}

func cmdSUNION(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SUNION' command"), false)
	}

	var sets []*data_structure.SimpleSet
	for _, key := range args {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].Union(sets[1:]...)

	return Encode(resultSet.Members(), false)
}

func cmdSUNIONSTORE(args []string) []byte {
	if len(args) < 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SUNIONSTORE' command"), false)
	}

	destKey := args[0]
	var sets []*data_structure.SimpleSet

	for _, key := range args[1:] {
		set, exist := setStore[key]
		if !exist {
			set = data_structure.NewSimpleSet(key)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].UnionStore(destKey, sets[1:]...)

	// Replace the old key with the new one
	delete(setStore, destKey)

	// Assign the new one
	setStore[destKey] = resultSet

	numbers := resultSet.Card()

	// Return the number of elems in the stored set
	return Encode(numbers, false)
}
