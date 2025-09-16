package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing/commands"
)

// REdis Serialization Protocol (RESP) is the simple and easy-to-parse protocol using by Redis clients and servers to communicate with each other.
// Unlike HTTP or HTTPS, the requests contain verbose meta-data (header fields) in it
//  --> unnecessary parsing
//	--> Increase latency a little bit

// Supported data types to represent for each type
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

const CRLF string = "\r\n"

var RespNil = []byte("$-1\r\n")

// +OK\r\n ==> "OK", 5
func readString(data []byte) (string, int, error) {
	return string(data[1 : len(data)-2]), len(data), nil
}

// Có thể sẽ phải implement lại kiểu khác

// :196\r\n ==> 196, 6
// :123\r\n => 123, 5
// :-99\r\n => -99, 6
// res: the length of the bytes array
// len(data) ~ pos: the last position of the array

func readInt(data []byte) (int64, int, error) {
	var val int64 = 0
	pos := 1
	var sign int64 = 1

	if data[pos] == '-' {
		sign = -1
		pos++
	}

	for data[pos] != '\r' {
		val = val*10 + int64(data[pos]-'0')
		pos++
	}

	return sign * val, pos + 2, nil
}

// Helper func
// $10\r\nHello\r\nYou\r\n ==> 10, 5
func readLen(data []byte) (int, int) {
	res, pos, _ := readInt(data)

	return int(res), pos

}

func readError(data []byte) (string, int, error) {
	return readString(data)
}

// $12\r\nHello\r\nWorld\r\n ==> "Hello\r\nWorld"
func readBulkString(data []byte) (string, int, error) {
	len, pos := readLen(data)
	s := string(data[pos:(pos + len)])

	return s, len + pos + 2, nil
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n => {"hello", "world"}
func readArray(data []byte) (interface{}, int, error) {
	length, pos := readLen(data)

	var res []interface{} = make([]interface{}, length)

	for i := range res {
		elem, delta, err := DecodeOne(data[pos:])

		if err != nil {
			return nil, 0, err
		}

		res[i] = elem
		pos += delta
	}

	return res, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case STRING:
		return readString(data)
	case INTEGER:
		return readInt(data)
	case ERROR:
		return readError(data)
	case BULK:
		return readBulkString(data)
	case ARRAY:
		return readArray(data)
	}
	return nil, 0, nil
}

// RESP format data => raw data
func Decode(data []byte) (interface{}, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}
	return []byte(fmt.Sprintf("*%d\r\n%s", len(sa), buf.Bytes()))
}

// raw data => RESP format data
func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, CRLF))
		}
		return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF))
	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case []string:
		return encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range value.([][]string) {
			buf.Write(encodeStringArray(sa))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([][]string)), buf.Bytes()))
	case []interface{}:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range value.([]interface{}) {
			buf.Write(Encode(x, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([]interface{})), buf.Bytes()))
	default:
		return RespNil
	}
}

func ParseCmd(data []byte) (*commands.Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]interface{})
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &commands.Command{Cmd: strings.ToUpper(tokens[0]), Args: tokens[1:]}
	return res, nil
}
