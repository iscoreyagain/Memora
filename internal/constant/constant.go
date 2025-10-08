package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespOk = []byte("+OK\r\n")
var RespZero = []byte(":0\r\n")
var RespOne = []byte(":1\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireFrequency = 100 * time.Millisecond
var ActiveExpireSampleSize = 20
var ActiveExpireThreshold = 0.1
var DefaultBPlusTreeDegree = 4

const BfDefaultInitCapacity = 100
const BfDefaultErrRate = 0.01

const CfDefaultInitCapacity = 128 // (1 << 7)
const CfDefaultBucketSize = 8

const SERVER_IDLE = 1
const SERVER_BUSY = 2
const SERVER_SHUTDOWN = 3

const LP_HEADER_SIZE = 6
const ENCODING_7BIT_UINT byte = 0    // 0xxxxxxx
const ENCODING_6BIT_STR byte = 0x80  // 10xxxxxx
const ENCODING_16BIT_INT byte = 0xF1 // 11110001
const ENCODING_32BIT_INT byte = 0xF3 // 11110011
const ENCODING_64BIT_INT byte = 0xF4 // 11110100
const ENCODING_12BIT_STR byte = 0xE0 // 1110xxxx
const ENCODING_32BIT_STR byte = 0xF0 // 11110000

// String length
