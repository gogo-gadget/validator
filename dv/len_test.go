package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type lenTest struct {
	name   string
	value  reflect.Value
	subTag string
}

func TestValidateLen(t *testing.T) {
	stringVal := "0123456789"
	mapVal := map[string]interface{}{"0": 0, "one": "1", "II": true}
	arrVal := []interface{}{"0", 1, false}

	values := []lenTest{
		{
			name:   "string",
			value:  reflect.ValueOf(stringVal),
			subTag: "len(10)",
		},
		{
			name:   "pointer",
			value:  reflect.ValueOf(&stringVal),
			subTag: "len(10)",
		},
		{
			name:   "map",
			value:  reflect.ValueOf(mapVal),
			subTag: "len(3)",
		},
		{
			name:   "array",
			value:  reflect.ValueOf(arrVal),
			subTag: "len(3)",
		},
	}

	for _, test := range values {
		t.Run(test.name, func(t *testing.T) {
			f := &cv.Field{
				Value: test.value,
			}

			err := ValidateLen(context.Background(), f, &cv.ValidationContext{SubTag: test.subTag})

			assert.NoError(t, err)
		})
	}
}

func TestValidateLen_failsForInvalidKinds(t *testing.T) {
	values := []lenTest{
		{
			name:   "int",
			value:  reflect.ValueOf(1),
			subTag: "len(1)",
		},
		{
			name:   "struct",
			value:  reflect.ValueOf(struct{}{}),
			subTag: "len(1)",
		},
	}

	for _, test := range values {
		t.Run(test.name, func(t *testing.T) {
			f := &cv.Field{
				Value: test.value,
			}

			err := ValidateLen(context.Background(), f, &cv.ValidationContext{SubTag: test.subTag})

			assert.Error(t, err)
		})
	}
}
