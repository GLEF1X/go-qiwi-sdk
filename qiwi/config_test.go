package qiwi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("SomeToken", "+1111111111")
	assert.NoError(t, err)
	assert.IsType(t, &Config{}, c)
}

func TestNewConfigAddPlusToPhoneNumber(t *testing.T) {
	c, err := NewConfig("SomeToken", "1111111111")
	assert.NoError(t, err)
	assert.Equal(t, "+1111111111", c.PhoneNumber)
}

func TestFailNewConfigIfTokenIsEmpty(t *testing.T) {
	_, err := NewConfig("", "+1111111111")
	if assert.Error(t, err) {
		assert.Equal(t, errors.Is(err, ConfigInputInvalidErr), true)
	}
}

func TestFailNewConfigDueToInvalidPhoneNumberFormat(t *testing.T) {
	_, err := NewConfig("SomeToken", "+")
	if assert.Error(t, err) {
		assert.Equal(t, errors.Is(err, ConfigInputInvalidErr), true)
	}
}
