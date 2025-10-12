package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/iscoreyagain/Memora/internal/data_structure"
)

func cmdCFEXISTS(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.EXISTS' command"), false)
	}

	key, item := args[0], args[1]
	cuckoo, exist := cuckooStore[key]
	if !exist {
		return Encode(0, false)
	}

	ok := cuckoo.Exist(item)
	if !ok {
		return Encode(0, false)
	}
	return Encode(1, false)
}

func cmdCFADD(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.ADD' command"), false)
	}

	key, item := args[0], args[1]
	cuckoo, exist := cuckooStore[key]
	if !exist {
		cuckoo = data_structure.CreateCuckooFilter(constant.CfDefaultErrRate, constant.CfDefaultInitCapacity)
		cuckooStore[key] = cuckoo
	}

	ok := cuckoo.Add(item)
	if !ok {
		return Encode(errors.New(fmt.Sprintf("(error) the cuckoo '%s' are currently full. Please create another one", key)), false)
	} else {
		cuckoo.Items++
		return Encode(1, false)
	}
}

func cmdCFADDNX(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.ADDNX' command"), false)
	}

	key, item := args[0], args[1]
	cuckoo, exist := cuckooStore[key]
	if !exist {
		cuckoo = data_structure.CreateCuckooFilter(constant.CfDefaultErrRate, constant.CfDefaultInitCapacity)
		cuckooStore[key] = cuckoo
	}

	if ok := cuckoo.Exist(item); ok {
		return Encode(0, false)
	}

	ok := cuckoo.Add(item)
	if !ok {
		return Encode(errors.New(fmt.Sprintf("(error) the cuckoo '%s' are currently full. Please create another one", key)), false)
	} else {
		cuckoo.Items++
		return Encode(1, false)
	}
}

func cmdCFMEXISTS(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.MEXISTS' command"), false)
	}
	key := args[0]
	cuckoo, exist := cuckooStore[key]

	res := make([]string, len(args)-1)
	if !exist {
		for i := range res {
			res[i] = "0"
		}

		return Encode(res, false)
	}

	for i := 1; i < len(args); i++ {
		if cuckoo.Exist(args[i]) {
			res[i-1] = "1"
		} else {
			res[i-1] = "0"
		}
	}
	return Encode(res, false)
}

func cmdCFRESERVE(args []string) []byte {
	if !(len(args) == 3 || len(args) == 5) {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.RESERVE' command"), false)
	}
	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("error rate must be a floating point number %s", args[1])), false)
	}
	capacity, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("capacity must be an integer number %s", args[2])), false)
	}
	_, exist := cuckooStore[key]
	if exist {
		return Encode(errors.New(fmt.Sprintf("Cuckoo filter with key '%s' already exist", key)), false)
	}
	cuckooStore[key] = data_structure.CreateCuckooFilter(errRate, int(capacity))
	return constant.RespOk
}

func cmdCFDEL(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CF.DEL' command"), false)
	}

	key, item := args[0], args[1]
	cuckoo, exist := cuckooStore[key]
	if !exist {
		return Encode(0, false)
	}

	ok := cuckoo.Remove(item)
	if !ok {
		return Encode(0, false)
	}
	return Encode(1, false)
}
