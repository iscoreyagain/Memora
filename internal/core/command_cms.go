package core

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/iscoreyagain/Memora/internal/data_structure"
)

func cmdCMSINITBYDIM(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYDIM' command"), false)
	}
	key := args[0]
	width, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("width must be a integer number %s", args[1])), false)
	}
	height, err := strconv.ParseUint(args[2], 10, 32)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("height must be a integer number %s", args[1])), false)
	}
	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exists"), false)
	}
	cmsStore[key] = data_structure.CreateCMS(uint64(width), uint64(height))
	return constant.RespOk
}

func cmdCMSINITBYPROB(args []string) []byte {
	if len(args) != 3 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'CMS.INITBYDIM' command"), false)
	}

	key := args[0]
	errRate, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("errRate must be a floating point number %s", args[1])), false)
	}
	if errRate >= 1 || errRate <= 0 {
		return Encode(errors.New("CMS: invalid overestimation value"), false)
	}
	probability, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return Encode(errors.New(fmt.Sprintf("probability must be a floating poit number %s", args[2])), false)
	}
	if probability >= 1 || probability <= 0 {
		return Encode(errors.New("CMS: invalid prob value"), false)
	}
	_, exist := cmsStore[key]
	if exist {
		return Encode(errors.New("CMS: key already exists"), false)
	}
	w, h := data_structure.CalcCMSDim(errRate, probability)
	cmsStore[key] = data_structure.CreateCMS(w, h)
	return constant.RespOk
}
