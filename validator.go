package validator

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Field struct {
	structField reflect.StructField
	value       reflect.Value
}

type CustomValidationFunc func(ctx context.Context, f Field) error

type CustomValidator struct {
	Validate CustomValidationFunc
	TagRegex *regexp.Regexp
	shouldFailIfNil bool
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

	// if the kind of the provided interface is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr{
		if iValue.IsNil() {
			// TODO fail validators that should fail on a nil ptr
			return nil
		}

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

func (v *Validator) RegisterValidationFunc(cv CustomValidator){
	if v.CustomValidators == nil{
		v.CustomValidators = map[string]CustomValidator{}
	}

	v.CustomValidators[cv.TagRegex.String()] = cv
}

func (v *Validator) validateStruct(ctx context.Context, structType reflect.Type, structValue reflect.Value) error {
	for i:=0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		field := Field{
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

func (v *Validator) validateField(ctx context.Context, field Field) error {
	// Validate Field if it contains a subTag matching a regex of any custom validator
	for _, cv := range v.CustomValidators {
		validatorTag := field.structField.Tag.Get("validator")
		subTags := strings.Split(validatorTag, ";")

		for _, subTag := range subTags{
			if cv.TagRegex.MatchString(subTag) {
				err := cv.Validate(ctx, field)
				if err != nil{
					return err
				}
			}
		}
	}

	fType := field.structField.Type
	fValue := field.value
	kind := fType.Kind()


	// if the kind of the field is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr{
		if fValue.IsNil() {
			// TODO fail validators that should fail on a nil ptr
			return nil
		}

		fType = fType.Elem()
		fValue = fValue.Elem()
		kind = fType.Kind()
	}

	// If the field itself is of kind struct validate the nested struct
	if kind == reflect.Struct {
		err := v.validateStruct(ctx, fType, fValue)
		if err != nil {
			return err
		}
	}

	return nil
}