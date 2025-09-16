package commands

import (
	"errors"
	"fmt"
	"strconv"
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
	cmsStore[key] = data_structure.CreateCMS(uint32(width), uint32(height))
	return constant.RespOk
}
