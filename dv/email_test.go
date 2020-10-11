package dv

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type CustomString string

func TestValidateEmail(t *testing.T) {
	email := "test@test.com"
	customString := CustomString(email)

	values := map[string]reflect.Value{
		"string": reflect.ValueOf(email),
		"pointer": reflect.ValueOf(&email),
		"custom type": reflect.ValueOf(customString),
	}

	for key, val := range values{
		t.Run(key, func(t *testing.T) {
			f := &cv.Field{
				Value: val,
			}

			err := ValidateEmail(context.Background(), f)

			assert.NoError(t, err)
		})
	}
}

func TestValidateEmail_failsForWrongType(t *testing.T){
	intVal := 1
	var interfaceVal interface{}
	interfaceVal = intVal

	values := map[string]reflect.Value{
		"int": reflect.ValueOf(intVal),
		"pointer to int": reflect.ValueOf(&intVal),
		"interface": reflect.ValueOf(interfaceVal),
	}

	for key, val := range values{
		t.Run(key, func(t *testing.T) {
			f := &cv.Field{
				Value: val,
			}

			err := ValidateEmail(context.Background(), f)

			assert.Error(t, err)
		})
	}
}

func TestValidateEmail_failsForInvalidEmail(t *testing.T) {
	values := map[string]reflect.Value{
		"int string": reflect.ValueOf("123456"),
		"phone number": reflect.ValueOf("+49123456789"),
		"postal code": reflect.ValueOf("12345"),
		"multiple dots": reflect.ValueOf("test@test..com"),
		"multiple ats": reflect.ValueOf("test@@test.com"),
		"whitespace before at": reflect.ValueOf("test @test.com"),
		"tab": reflect.ValueOf("test	@test.com"),
		"whitespace after at": reflect.ValueOf("test@ test.com"),
	}

	for key, val := range values{
		t.Run(key, func(t *testing.T) {
			f := &cv.Field{
				Value: val,
			}

			err := ValidateEmail(context.Background(), f)

			assert.Error(t, err)
		})
	}
}