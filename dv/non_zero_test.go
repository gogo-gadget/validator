package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

func TestValidateNonZero(t *testing.T) {
	f := &cv.Field{
		Value: reflect.ValueOf("non-zero"),
	}

	err := ValidateNonZero(context.Background(), f)

	assert.NoError(t, err)
}

func TestValidateNonZero_failsForZeroValue(t *testing.T) {
	f := &cv.Field{
		Value: reflect.ValueOf(""),
	}

	err := ValidateNonZero(context.Background(), f)

	assert.Error(t, err)
}
