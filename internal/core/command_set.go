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
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
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
	if len(args) <= 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SINTER' command"), false)
	}

	var sets []*data_structure.SimpleSet

	for _, key := range args {
		set, exist := setStore[key]
		if !exist {
			return Encode(nil, false)
		}
		sets = append(sets, set)
	}

	resultSet := sets[0].Intersection(sets[1:]...)

	members := make([]string, 0, len(resultSet.Members()))
	for _, m := range resultSet.Members() {
		members = append(members, m)
	}

	return Encode(members, false)
}
