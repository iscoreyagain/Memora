package data_structure

import (
	"testing"

	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/stretchr/testify/assert"
)

func TestBacklenFunctions(t *testing.T) {
	buf := make([]byte, 5) // enough capacity

	n := EncodeBacklen(buf, 1045241)

	result := DecodeBacklen(buf[:n])

	expected := uint64(1045241)
	assert.Equal(t, expected, result, "Decoded value should match the encoded value")
}

func TestEncodingInteger(t *testing.T) {
	b, n := encodingInteger(-2)
	assert.Equal(t, int64(3), n)
	assert.Equal(t, []byte{constant.ENCODING_16BIT_INT, 0xFE, 0xFF}, b)
}
