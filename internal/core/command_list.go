package core

import (
	"errors"

	"github.com/iscoreyagain/Memora/internal/data_structure"
)

func cmdRPUSH(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'RPUSH' command"), false)
	}

	key := args[0]
	list, exist := listStore[key]
	if !exist {
		list = data_structure.NewList(key)
		listStore[key] = list
	}
	list.RPush(convertToInterfaces(args[1:])...)

	return Encode(list.Len(), false)
}

func cmdLLEN(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'LLEN' command"), false)
	}

	key := args[0]
	list, exist := listStore[key]
	if !exist {
		return Encode(0, false)
	}

	count := list.Len()
	return Encode(count, false)
}

func cmdLPUSH(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'RPUSH' command"), false)
	}

	key := args[0]
	list, exist := listStore[key]
	if !exist {
		list = data_structure.NewList(key)
		listStore[key] = list
	}
	list.LPush(convertToInterfaces(args[1:])...)

	return Encode(list.Len(), false)
}

// helper functions
func convertToInterfaces(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, s := range strs {
		res[i] = s
	}
	return res
}
