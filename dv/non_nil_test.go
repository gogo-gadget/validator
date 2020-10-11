package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type simpleNonNilStruct struct {
	Simple *simpleNonNilStruct `validator:"non-nil"`
}

func TestValidateNonNil(t *testing.T) {
	f := &cv.Field{
		Value: reflect.ValueOf(&simpleNonNilStruct{}),
	}

	err := ValidateNonNil(context.Background(), f, &cv.ValidationContext{})

	assert.NoError(t, err)
}

func TestValidateNonNil_failsForNilValue(t *testing.T) {
	nilValue := reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
	f := &cv.Field{
		Value: nilValue,
	}

	err := ValidateNonNil(context.Background(), f, &cv.ValidationContext{})

	assert.Error(t, err)
}
