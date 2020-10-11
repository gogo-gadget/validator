package validator

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gogo-gadget/validator/pkg/cv"
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
	kind := iValue.Kind()

	// if the kind of the provided interface is interface or pointer use its underlying element instead
	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if iValue.IsNil() {
			// fail validators that should fail on a nil ptr
			iType = getUnderlyingType(iType)
			kind = iType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(iType, nil)
			if err != nil {
				return err
			}

			return nil
		}

		iValue = iValue.Elem()
		iType = iValue.Type()
		kind = iValue.Kind()
	}

	if kind != reflect.Struct {
		return fmt.Errorf("validation of kind %v is not supported", kind)
	}

	err := v.validateStruct(ctx, iValue, nil)
	if err != nil {
		return err
	}

	return nil
}

// Should only be used on reflect.Values of kind struct
func (v *Validator) validateStruct(ctx context.Context, structValue reflect.Value, parent *cv.Field) error {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		field := &cv.Field{
			Parent:      parent,
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

func (v *Validator) validateField(ctx context.Context, field *cv.Field) error {
	// Validate Field if it contains a subTag matching a regex of any custom validator
	validatorTag := field.StructField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, customValidator := range v.CustomValidators {
		for _, subTag := range subTags {
			if customValidator.TagRegex.MatchString(subTag) {
				validationCtx := &cv.ValidationContext{
					Tag:    validatorTag,
					SubTag: subTag,
				}
				err := customValidator.Validate(ctx, field, validationCtx)
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
	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if fValue.IsNil() {
			// fail validators that should fail on a nil ptr
			fType = getUnderlyingType(fType)
			kind = fType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(fType, field)
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
	err := v.validateStruct(ctx, fValue, field)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateStructNilValidations(structType reflect.Type, parent *cv.Field) error {
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		field := &cv.Field{
			Parent:      parent,
			StructField: structField,
		}

		err := v.validateFieldNilValidations(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateFieldNilValidations(field *cv.Field) error {
	structField := field.StructField
	validatorTag := structField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, customValidator := range v.CustomValidators {
		for _, subTag := range subTags {
			if customValidator.TagRegex.MatchString(subTag) && customValidator.Config.ShouldFailIfFieldOfNilPtr {
				fullFieldName := getFullFieldName(field)
				return fmt.Errorf("validation failed since validator for regex: %v failed on nil value for field: %v", customValidator.TagRegex.String(), fullFieldName)
			}
		}
	}

	fType := getUnderlyingType(structField.Type)
	kind := fType.Kind()

	// If the field is not of kind struct there is nothing to be validated anymore
	if kind != reflect.Struct {
		return nil
	}

	// If the field itself is of kind struct validate the nested struct
	err := v.validateStructNilValidations(fType, field)
	if err != nil {
		return err
	}

	return nil
}

func getFullFieldName(field *cv.Field) string {
	if field == nil {
		return ""
	}

	fullFieldName := field.StructField.Name
	parent := field.Parent
	for parent != nil {
		fullFieldName = fmt.Sprintf("%v.%v", parent.StructField.Name, fullFieldName)
		parent = parent.Parent
	}

	return fullFieldName
}

func getUnderlyingType(rType reflect.Type) reflect.Type {
	kind := rType.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		rType = rType.Elem()
		kind = rType.Kind()
	}

	return rType
}
