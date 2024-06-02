package spine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperations(t *testing.T) {
	operations := NewOperations(true, true, false, false)
	assert.NotNil(t, operations)

	text := operations.String()
	assert.NotEqual(t, 0, len(text))

	data := operations.Information()
	assert.NotNil(t, data)

	result := operations.Read()
	assert.True(t, result)
	result = operations.ReadPartial()
	assert.True(t, result)
	result = operations.Write()
	assert.False(t, result)
	result = operations.WritePartial()
	assert.False(t, result)

	operations2 := NewOperations(true, false, true, true)
	assert.NotNil(t, operations2)

	text = operations2.String()
	assert.NotEqual(t, 0, len(text))

	data = operations2.Information()
	assert.NotNil(t, data)

	result = operations2.Read()
	assert.True(t, result)
	result = operations2.ReadPartial()
	assert.False(t, result)
	result = operations2.Write()
	assert.True(t, result)
	result = operations2.WritePartial()
	assert.True(t, result)

	operations3 := NewOperations(false, false, false, false)
	assert.NotNil(t, operations3)

	text = operations3.String()
	assert.NotEqual(t, 0, len(text))

	data = operations3.Information()
	assert.NotNil(t, data)

	result = operations3.Read()
	assert.False(t, result)
	result = operations3.ReadPartial()
	assert.False(t, result)
	result = operations3.Write()
	assert.False(t, result)
	result = operations3.WritePartial()
	assert.False(t, result)
}
