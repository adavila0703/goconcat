package goconcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetJSONSettings(t *testing.T) {
	assert := assert.New(t)

	options := NewOptions()

	err := options.SetJSONOptions("optionsExample.json")
	assert.NoError(err)

	assert.Equal(".", options.RootPath)
}

func TestSetMockeryDestination(t *testing.T) {
	assert := assert.New(t)

	options := NewOptions()

	err := options.SetJSONOptions("optionsExample.json")
	assert.NoError(err)

	assert.Equal(".", options.RootPath)
	options.SetMockeryDestination(true)

	assert.Equal(true, options.MockeryDestination)
}
