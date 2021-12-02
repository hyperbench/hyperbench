package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNet(t *testing.T) {
	s := Bytes2Hex(nil)
	assert.Equal(t, s, "")
	bs := Hex2Bytes("nil")
	assert.NotNil(t, bs)
}
