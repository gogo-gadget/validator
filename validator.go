package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
)

type field struct {
	structField reflect.StructField
	value       reflect.Value
}

type CustomValidationFunc func(ctx context.Context, f field) error

type CustomValidator struct {
	Validate CustomValidationFunc
	TagRegex regexp.Regexp
}

type Validator struct {
	CustomValidators map[string]CustomValidator
}

func NewValidator() *Validator{
	return &Validator{}
}

func (v *Validator) Validate(ctx context.Context, i interface{}) error{
	iType := reflect.TypeOf(i)
	iValue := reflect.ValueOf(i)

	kind := iType.Kind()

	if kind == reflect.Interface || kind == reflect.Ptr{
		iType = iType.Elem()
		iValue = iValue.Elem()
		kind = iType.Kind()
	}

	if kind != reflect.Struct{
		return fmt.Errorf("validation of kind %v is not supported", kind)
	}

	err := v.validateStruct(ctx, iType, iValue)
	if err != nil{
		return err
	}

	return nil
}

func (v *Validator) RegisterValidationFunc(tagRegex regexp.Regexp, cv CustomValidator){
	v.CustomValidators[tagRegex.String()] = cv
}

func (v *Validator) validateStruct(ctx context.Context, structType reflect.Type, structValue reflect.Value) error {
	for i:=0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		field := field{
			structField: structField,
			value:       fieldValue,
		}
		
		err := v.validateField(ctx, field)
		if err != nil{
			return err
		}
	}

	return nil
}

func (v *Validator) validateField(ctx context.Context, field field) error {
	for _, cv := range v.CustomValidators {
		err := cv.Validate(ctx, field)
		if err != nil{
			return err
		}
	}

	return nil
}