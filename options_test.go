package goconcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetJSONSettings(t *testing.T) {
	assert := assert.New(t)

	options := NewOptions()

	err := options.SetJSONSettings("optionsExample.json")
	assert.NoError(err)

	assert.Equal(".", options.RootPath)
}
