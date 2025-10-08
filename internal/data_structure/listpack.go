package data_structure

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/iscoreyagain/Memora/internal/constant"
)

const EOF = 0xFF

// Listpack is memory-efficient data structure used by Redis to store a small amounts of element
// in a LIST. This creation due to overhead of allocating heap memory for struct and other fields
// as well as replacing `Ziplist` in Redis's internal implementation - complex + bugs prone
// Understanding the efficient idea of Ziplist, Listpack is conceptually a compacted byte array contains:

// ******************** Header (6 bytes) ********************
// 	+ tot-bytes (4 bytes) tell us the total length of the listpack (including the header itself + the terminator byte)
//  + num-elements (2 bytes) give us the count of elements that this listpack currently stored

// - entry-len: the total length of the entry (including itself) for jumping
// - encoding: the internal-used flag tell us what the data type of entry stored
// - data/content: the actual content
// - 0xFF: the end marker after the last entry

//	[entry-len][data]
//  in `entry-len`: [encoding][extra length bytes]

type ListPack struct {
	data []byte
}

func SetTotBytesLp(lp *ListPack, b uint32) {
	binary.LittleEndian.PutUint32(lp.data[0:4], b)
}

func SetNumElemsLp(lp *ListPack, b uint16) {
	binary.LittleEndian.PutUint16(lp.data[4:6], b)
}

func NewListPack() *ListPack {
	data := make([]byte, constant.LP_HEADER_SIZE+1)
	data[constant.LP_HEADER_SIZE] = 0xFF
	lp := &ListPack{data: data}
	SetTotBytesLp(lp, constant.LP_HEADER_SIZE+1)
	SetNumElemsLp(lp, 0)
	return lp
}

func (lp *ListPack) Skip(pos uint64) uint64 {
	entryLen := GetCurrentEncodedSize(lp.data[pos:])
	entryLen += encodeBacklen(nil, entryLen)

	return pos + entryLen
}

func (lp *ListPack) Next(pos uint64) *uint64 {
	pos += lp.Skip(pos)
	if lp.data[pos] == EOF {
		return nil
	}
	return &pos
}

func (lp *ListPack) Prev(pos uint64) *uint64 {
	if pos == 0 || lp.data[pos-1] == EOF {
		return nil
	}

	prevLen := decodeBacklen(lp.data[:pos])
	total := prevLen + encodeBacklen(nil, prevLen)

	newPos := pos - total

	return &newPos
}

func GetTotBytes(lp *ListPack) int32 {
	return int32(binary.LittleEndian.Uint32(lp.data[:4]))
}

func GetNumElems(lp *ListPack) int16 {
	return int16(binary.LittleEndian.Uint16(lp.data[4:6]))
}

// LPUSH
func (lp *ListPack) LPush(members ...interface{}) int {
	added := 0
	var arr []byte
	for _, m := range members {
		entry, err := EncodeEntry(m)
		if err != nil {
			continue
		}
		arr = append(arr, entry...)
		added++
	}
	lp.data = append(arr, lp.data...)

	if len(lp.data) == 0 || lp.data[len(lp.data)-1] != 0xFF {
		lp.data = append(lp.data, 0xFF)
	}

	return added
}

// RPUSH
func (lp *ListPack) RPush(members ...interface{}) int {
	added := 0

	if len(lp.data) > 0 && lp.data[len(lp.data)-1] == 0xFF {
		lp.data = lp.data[:len(lp.data)-1]
	}

	for _, m := range members {
		entry, err := EncodeEntry(m)
		if err != nil {
			return added
		}
		lp.data = append(lp.data, entry...)
		added++
	}

	lp.data = append(lp.data, 0xFF)
	return added
}

func EncodeEntry(member interface{}) ([]byte, error) {
	var encoding byte
	var content []byte
	var entryLen int
	var err error

	switch v := member.(type) {
	case string:
		encoding, content, entryLen, err = encodeString(v)
		if err != nil {
			return nil, err
		}
	case int8, int16, int32, int64, int:
		encoding, content, entryLen, err = encodeInteger(member)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported type: %T", member)
	}

	// 1. Size prefix - tell us how far we should jump to move to the next entry
	// It's usually 1 OR 5 bytes
	sizePrefix := encodeLen(entryLen)

	// 2. Assemble to get the encoded entry
	totalSize := len(content) + len(sizePrefix) + 1 // 1-byte of encoding
	entry := make([]byte, totalSize)
	entry = append(entry, sizePrefix...)
	entry = append(entry, encoding)
	entry = append(entry, content...)

	return entry, nil
}

func encodeLen(length int) []byte {
	// One-byte length (for small entries)
	// length = 25 ---> [0x19]
	if length <= 127 {
		return []byte{byte(length)}
	}

	// 5-byte length (for large entries)
	// length = 300 ----> [0xFE   0x00   0x00   0x01   0x2C] ----> Little Endian ----> [0xFE   0x2C   0x01   0x00   0x00]
	b := make([]byte, 5)
	b[0] = 0xFE
	binary.LittleEndian.PutUint32(b[1:], uint32(length))

	return b
}

func encodeString(member string) (byte, []byte, int, error) {
	var encoding byte
	var content []byte
	var entryLen int

	length := len(member)
	if length < 64 {
		encoding = byte(length) & 0x3F // 0011 1111 - for .i.e 10010110 & 00111111 -> 00010110
		content = []byte(member)
		entryLen = 1 + length
	} else if length < 16384 { // for the case length < 16384
		encoding = byte(length>>8) & 0x3F
		encoding |= 0x40 // 0100 0000
		content = append([]byte{byte(length & 0xFF)}, []byte(member)...)
		entryLen = 1 + len(content)
	} else if length <= (1 << 32) {
		// 0x80 -> 10|00 0000
		encoding = byte(0x80)
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, uint32(length))
		content = append(lenBytes, []byte(member)...)
		entryLen = 1 + len(content)
	} else {
		return 0, nil, 0, fmt.Errorf("string length exceeds 32-bit limit")
	}
	return encoding, content, entryLen, nil
}

func encodeInteger(member interface{}) (byte, []byte, int, error) {
	var encoding byte
	var content []byte
	var entryLen int

	switch v := member.(type) {
	case int8: // 1 byte
		entryLen = 1 + 1
		encoding = 0xC0
		content = []byte{byte(v)}

	case int16: // 2 byte
		entryLen = 1 + 2
		encoding = 0xD0
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v))
		content = b

	// Go doesn't have supported int24 (3-byte int) yet so for educational purpose, we will temporarily ignore it
	case int32:
		entryLen = 1 + 4
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		encoding = 0xF0
		content = b

	case int64:
		encoding = 0xF1
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v))
		entryLen = 1 + 8
		content = b
	default:
		return 0, nil, 0, fmt.Errorf("Unsupported integer type: %T", member)
	}

	return encoding, content, entryLen, nil
}

func Get6BitStrLen(bytes []byte) int {
	return int(bytes[0] & 0x3F)
}

func Get12BitStrLen(bytes []byte) uint16 {
	return (uint16(bytes[0]&0x0F) << 8) | uint16(bytes[1])
}

func Get32BitStrLen(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes[1:5])
}

// *** Friendly reminder for anyone to understand deep dive into the format that - this function is
// currently calculate the length of the current element WITHOUT thr `backlen` field
func GetCurrentEncodedSize(bytes []byte) uint64 {
	if bytes[0] == EOF {
		return 1
	}
	if (bytes[0] & 0x80) == constant.ENCODING_7BIT_UINT {
		return 1
	}
	// To be more specific, since we only care the first two bits to determine this type is 6-bit string
	// We would only need to have a mask corresponding to that 2 position: 11000000
	if (bytes[0] & 0xC0) == constant.ENCODING_6BIT_STR {
		return uint64(1 + Get6BitStrLen(bytes))
	}
	if (bytes[0] & 0xFF) == constant.ENCODING_16BIT_INT {
		return 3
	}
	if (bytes[0] & 0xF0) == constant.ENCODING_12BIT_STR {
		return uint64(2 + Get12BitStrLen(bytes))
	}
	if (bytes[0] & 0xFF) == constant.ENCODING_32BIT_INT {
		return 5
	}
	if (bytes[0] & 0xFF) == constant.ENCODING_64BIT_INT {
		return 9
	}
	if (bytes[0] & 0xFF) == constant.ENCODING_32BIT_STR {
		return uint64(5 + Get32BitStrLen(bytes))
	}
	return 0
}

// This is the last thing of a LP's element structure (encoding type + element data + backlen) to be added. Thanks to
// this genius design allowing us to backward traversal (traverse from the right to the left)
func encodeBacklen(buf []byte, length uint64) uint64 {
	if length < 127 {
		if buf != nil {
			buf[0] = byte(length)
		}
		return 1
	} else if length < 16383 {
		if buf != nil {
			buf[0] = byte(length >> 7)
			buf[1] = byte(length&0xFF) | 128
		}
		return 2
	} else if length < 2097151 {
		if buf != nil {
			buf[0] = byte(length >> 14)
			buf[1] = byte((length>>7)&0xFF) | 128
			buf[2] = byte(length&0xFF) | 128
		}
		return 3
	} else if length < 268435455 {
		if buf != nil {
			buf[0] = byte(length >> 21)
			buf[1] = byte((length>>14)&0xFF) | 128
			buf[2] = byte((length>>7)&0xFF) | 128
			buf[3] = byte(length&0xFF) | 128
		}
		return 4
	} else {
		if buf != nil {
			buf[0] = byte(length >> 28)
			buf[1] = byte((length>>21)&0xFF) | 128
			buf[2] = byte((length>>14)&0xFF) | 128
			buf[3] = byte((length>>7)&0xFF) | 128
			buf[4] = byte(length&0xFF) | 128
		}
		return 5
	}
}

func decodeBacklen(bytes []byte) uint64 {
	var value uint64
	var shift uint64

	for i := len(bytes) - 1; i >= 0; i-- {
		value |= uint64(bytes[i]&0x7F) << shift
		if bytes[i]&0x80 == 0 {
			break
		}
		shift += 7
		if shift > 28 {
			return math.MaxUint64
		}
	}

	return value
}
