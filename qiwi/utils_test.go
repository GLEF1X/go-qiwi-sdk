package qiwi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowerFirstLetter(t *testing.T) {
	assert.Equal(t, "hello", lowerFirstLetter("Hello"))
}

func TestLowerFirstLetterReturnIfEmptyStringInput(t *testing.T) {
	assert.Equal(t, "", lowerFirstLetter(""))
}
