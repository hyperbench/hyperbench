package vm

import (
	"testing"

	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/stretchr/testify/assert"
)

func TestNewVM(t *testing.T) {
	t.Skip()
	vm, err := NewVM("lua", base.ConfigBase{})
	assert.Nil(t, vm)
	assert.Error(t, err)

	vm, err = NewVM("", base.ConfigBase{})
	assert.NotNil(t, vm)
	assert.NoError(t, err)
}
