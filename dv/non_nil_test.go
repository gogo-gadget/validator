package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type simpleStruct struct {
	Simple *simpleStruct `validator:"non-nil"`
}

func TestValidateNonNil(t *testing.T) {
	f := &cv.Field{
		StructField: reflect.StructField{},
		Value:       reflect.ValueOf(&simpleStruct{}),
	}

	err := ValidateNonNil(context.Background(), f)

	assert.NoError(t, err)
}

func TestValidateNonNil_failsForNilValue(t *testing.T) {
	nilValue := reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
	f := &cv.Field{
		StructField: reflect.StructField{},
		Value:       nilValue,
	}

	err := ValidateNonNil(context.Background(), f)

	assert.Error(t, err)
}
