package validator

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ValidEmail = "test@test.com"
const InvalidEmail = "1234"

type And struct {
	Field string `validator:"len(13) && email"`
}

type And2 struct {
	Field string `validator:"email && len(13)"`
}

type Or struct {
	Field string `validator:"len(13) || email"`
}

type Or2 struct {
	Field string `validator:"email || len(13)"`
}

type Negate struct {
	Field string `validator:"!email"`
}

type Negate2 struct {
	Field string `validator:"!len(13)"`
}

type MultiOperators struct {
	Field string `validator:"required && email && len(13)"`
}

func TestValidator_Validate_logicalOperators(t *testing.T) {
	validator := NewValidator()
	invalidEmail := "1234567890123"

	values := []interface{}{
		And{Field: ValidEmail},
		And2{Field: ValidEmail},
		Or{Field: invalidEmail},
		Or2{Field: invalidEmail},
		Negate{Field: invalidEmail},
		Negate2{Field: "1234"},
		MultiOperators{Field: ValidEmail},
	}

	for _, val := range values {
		t.Run(reflect.TypeOf(val).Name(), func(t *testing.T) {
			err := validator.Validate(context.Background(), val)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_Validate_logicalOperatorsFailsCorrectly(t *testing.T) {
	validator := NewValidator()
	invalidEmail := "1234567890123"

	values := []interface{}{
		And{Field: invalidEmail},
		And2{Field: invalidEmail},
		Or{Field: ""},
		Or2{Field: ""},
		Negate{Field: ValidEmail},
		Negate2{Field: invalidEmail},
	}

	for _, val := range values {
		t.Run(reflect.TypeOf(val).Name(), func(t *testing.T) {
			err := validator.Validate(context.Background(), val)
			assert.Error(t, err)
		})
	}
}

type FrontAndBrace struct {
	Field string `validator:"(len(13) && email) && non-nil"`
}

type BackAndBrace struct {
	Field string `validator:"len(13) && (email && non-nil)"`
}

type FrontOrBrace struct {
	Field string `validator:"(email || required) && non-nil"`
}

type BackOrBrace struct {
	Field string `validator:"email || (required && non-nil)"`
}

func TestValidator_Validate_braces(t *testing.T) {
	validator := NewValidator()

	values := []interface{}{
		FrontAndBrace{Field: ValidEmail},
		BackAndBrace{Field: ValidEmail},
		FrontOrBrace{Field: InvalidEmail},
		BackOrBrace{Field: InvalidEmail},
	}

	for _, val := range values {
		t.Run(reflect.TypeOf(val).Name(), func(t *testing.T) {
			err := validator.Validate(context.Background(), val)
			assert.NoError(t, err)
		})
	}
}

type WhiteSpaceStruct struct {
	Field string `validator:"  email&&  required		"`
}

func TestValidator_Validate_ignoresWhitespace(t *testing.T) {
	validator := NewValidator()

	val := WhiteSpaceStruct{
		Field: ValidEmail,
	}

	err := validator.Validate(context.Background(), val)

	assert.NoError(t, err)
}

type IfStruct struct {
	Field string `validator:"if(email)then(len(13))"`
}

type IfElifStruct struct {
	Field string `validator:"if(!email)then(len(4))elif(email)then(len(13))"`
}

type IfElseStruct struct {
	Field string `validator:"if(email)then(len(13))"`
}

type IfElifElseStruct struct {
	Field string `validator:"if(len(13))then(email)elif(len(4))then(!email)else(!non-zero)"`
}

type IfAndStruct struct {
	Field string `validator:"if(email)then(len(13)) && non-zero"`
}

type IfOrStruct struct {
	Field string `validator:"if(email)then(len(4)) || non-zero"`
}

func TestValidator_Validate_ifStatement(t *testing.T) {
	validator := NewValidator()

	values := []interface{}{
		IfStruct{Field: ValidEmail},
		IfStruct{Field: InvalidEmail},
		IfElifStruct{Field: ValidEmail},
		IfElifStruct{Field: InvalidEmail},
		IfElseStruct{Field: ValidEmail},
		IfElseStruct{Field: InvalidEmail},
		IfAndStruct{Field: ValidEmail},
		IfAndStruct{Field: InvalidEmail},
		IfOrStruct{Field: ValidEmail},
		IfOrStruct{Field: InvalidEmail},
		IfElifElseStruct{Field: ValidEmail},
		IfElifElseStruct{Field: InvalidEmail},
		IfElifElseStruct{Field: ""},
	}

	for idx, val := range values {
		t.Run(fmt.Sprintf("subtest: %v", idx), func(t *testing.T) {
			err := validator.Validate(context.Background(), val)
			assert.NoError(t, err)
		})
	}
}
