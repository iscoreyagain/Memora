package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/iscoreyagain/Memora/internal/data_structure"
)

func cmdRPUSH(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'RPUSH' command"), false)
	}

	key := args[0]
	list, exist := listStore[key]
	if !exist {

	}
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
