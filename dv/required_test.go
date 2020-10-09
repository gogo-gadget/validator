package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

func TestRequired(t *testing.T) {
	f := &cv.Field{
		Value: reflect.ValueOf("required value"),
	}
	err := ValidateRequired(context.Background(), f)

	assert.NoError(t, err)
}

func TestRequired_failsForNilValue(t *testing.T) {
	nilValue := reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

	f := &cv.Field{
		Value: nilValue,
	}

	err := ValidateRequired(context.Background(), f)

	assert.Error(t, err)
}

func TestRequired_failsForZeroValue(t *testing.T) {
	f := &cv.Field{
		Value: reflect.ValueOf(""),
	}

	err := ValidateRequired(context.Background(), f)

	assert.Error(t, err)
}
