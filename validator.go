package validator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	cv "github.com/gogo-gadget/validator/pkg/custom-validator"
)

type Validator struct {
	CustomValidators map[string]*cv.CustomValidator
}

func NewValidator() *Validator {
	v := &Validator{}

	v.RegisterDefaultCustomValidators()

	return v
}

func (v *Validator) RegisterCustomValidator(customValidator *cv.CustomValidator) {
	if v.CustomValidators == nil {
		v.CustomValidators = map[string]*cv.CustomValidator{}
	}

	v.CustomValidators[customValidator.ID] = customValidator
}

func (v *Validator) Validate(ctx context.Context, i interface{}) error {
	iValue := reflect.ValueOf(i)
	iType := iValue.Type()
	kind := iType.Kind()

	// if the kind of the provided interface is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr {
		if iValue.IsNil() {
			// fail validators that should fail on a nil ptr
			iType = iType.Elem()
			kind = iType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(iType)
			if err != nil {
				return err
			}

			return nil
		}

		iValue = iValue.Elem()
		iType = iValue.Type()
		kind = iType.Kind()
	}

	if kind != reflect.Struct {
		return fmt.Errorf("validation of kind %v is not supported", kind)
	}

	err := v.validateStruct(ctx, iValue)
	if err != nil {
		return err
	}

	return nil
}

// Should only be used on reflect.Values of kind struct
func (v *Validator) validateStruct(ctx context.Context, structValue reflect.Value) error {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		field := cv.Field{
			StructField: structField,
			Value:       fieldValue,
		}

		err := v.validateField(ctx, field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateField(ctx context.Context, field cv.Field) error {
	// Validate Field if it contains a subTag matching a regex of any custom validator
	validatorTag := field.StructField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, customValidator := range v.CustomValidators {
		for _, subTag := range subTags {
			if customValidator.TagRegex.MatchString(subTag) {
				err := customValidator.Validate(ctx, field)
				if err != nil {
					return err
				}
			}
		}
	}

	fValue := field.Value
	fType := fValue.Type()
	kind := fType.Kind()

	// if the kind of the field is interface or pointer use its underlying element instead
	if kind == reflect.Interface || kind == reflect.Ptr {
		if fValue.IsNil() {
			// fail validators that should fail on a nil ptr
			fType = fType.Elem()
			kind = fType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(fType)
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

func (v *Validator) validateStructNilValidations(structType reflect.Type) error {
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		err := v.validateFieldNilValidations(structField)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateFieldNilValidations(structField reflect.StructField) error {
	validatorTag := structField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, customValidator := range v.CustomValidators {
		for _, subTag := range subTags {
			if customValidator.TagRegex.MatchString(subTag) && customValidator.Config.ShouldFailIfFieldOfNilPtr {
				return fmt.Errorf("validation failed since validator for regex: %v failed", customValidator.TagRegex.String())
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
	err := v.validateStructNilValidations(fType)
	if err != nil {
		return err
	}

	return nil
}
