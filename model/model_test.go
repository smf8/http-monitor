package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	foo, err := NewUser("Foo", "Bar")
	assert.NoError(t, err, "Error creating user instance")
	assert.False(t, foo.ValidatePassword("Bar"), "Error validating password")
}
