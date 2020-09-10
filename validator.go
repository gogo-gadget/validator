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
	Validate        CustomValidationFunc
	TagRegex        *regexp.Regexp
	ShouldFailIfNil bool
}

type Validator struct {
	CustomValidators map[string]CustomValidator
}

func NewValidator() *Validator{
	return &Validator{}
}

func (v *Validator) Validate(ctx context.Context, i interface{}) error{
	iValue := reflect.ValueOf(i)
	iType := iValue.Type()
	kind := iType.Kind()

	// if the kind of the provided interface is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr{
		if iValue.IsNil() {
			// TODO fail validators that should fail on a nil ptr
			// fail validators that should fail on a nil ptr
			iType = iType.Elem()
			kind = iType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.failStructNilValidations(iType)
			if err != nil {
				return err
			}

			return nil
		}

		iValue = iValue.Elem()
		iType = iValue.Type()
		kind = iType.Kind()
	}

	if kind != reflect.Struct{
		return fmt.Errorf("validation of kind %v is not supported", kind)
	}

	err := v.validateStruct(ctx, iValue)
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


// Should only be used on reflect.Values of kind struct
func (v *Validator) validateStruct(ctx context.Context, structValue reflect.Value) error {
	structType := structValue.Type()
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
	validatorTag := field.structField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, cv := range v.CustomValidators {
		for _, subTag := range subTags{
			if cv.TagRegex.MatchString(subTag) {
				err := cv.Validate(ctx, field)
				if err != nil{
					return err
				}
			}
		}
	}

	fValue := field.value
	fType := fValue.Type()
	kind := fType.Kind()


	// if the kind of the field is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr{
		if fValue.IsNil() {
			// fail validators that should fail on a nil ptr
			fType = fType.Elem()
			kind = fType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.failStructNilValidations(fType)
			if err != nil {
				return err
			}

			return nil
		}

		fValue = fValue.Elem()
		fType = fValue.Type()
		kind = fType.Kind()
	}

	// If the field is not of kind struct there is nothing to be validated anymore
	if kind != reflect.Struct {
		return nil
	}

	// If the field itself is of kind struct validate the nested struct
	err := v.validateStruct(ctx, fValue)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) failStructNilValidations(structType reflect.Type) error {
	for i:=0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		err := v.failStructFieldNilValidations(structField)
		if err != nil{
			return err
		}
	}

	return nil
}

func (v *Validator) failStructFieldNilValidations(structField reflect.StructField) error {
	validatorTag := structField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, cv := range v.CustomValidators{
		for _, subTag := range subTags {
			if cv.TagRegex.MatchString(subTag) && cv.ShouldFailIfNil {
				return fmt.Errorf("validation failed since validator for regex: %v failed", cv.TagRegex.String())
			}
		}
	}

	fType := structField.Type
	kind := fType.Kind()
	if kind == reflect.Interface || kind == reflect.Ptr {
		fType = fType.Elem()
		kind = fType.Kind()
	}

	// If the field is not of kind struct there is nothing to be validated anymore
	if kind != reflect.Struct {
		return nil
	}

	// If the field itself is of kind struct validate the nested struct
	err := v.failStructNilValidations(fType)
	if err != nil {
		return err
	}

	return nil
}