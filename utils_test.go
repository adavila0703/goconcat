package goconcat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToType(t *testing.T) {
	assert := assert.New(t)

	type testString string

	mockSlice := []string{"1", "2", "3"}

	respSlice := stringToType[testString](mockSlice)

	assert.Equal([]testString{"1", "2", "3"}, respSlice)
}

func TestPopFromSlice(t *testing.T) {
	assert := assert.New(t)

	mockSlice := []int{1, 2, 3}

	num := popFromSlice(mockSlice, 2)

	assert.Equal(3, num)
}
