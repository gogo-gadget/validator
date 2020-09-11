package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	cv "github.com/gogo-gadget/validator/pkg/custom-validator"
)

func NonNil() *cv.CustomValidator {
	nonNilString := "non-nil"
	nonNilRegexp := regexp.MustCompile(nonNilString)

	customValidator := cv.NewCustomValidator("non-nil", nonNilRegexp, ValidateNonNil, cv.NewCustomValidatorConfig().WithNilValidation(true))
	return customValidator
}

func ValidateNonNil(ctx context.Context, f cv.Field) error {
	kind := f.Value.Kind()

	switch kind {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// the fields value can only be nil if it is an interface, pointer, map, slice, chan, func or unsafe pointer
		if f.Value.IsNil() {
			return fmt.Errorf("NonNil field %v is nil", f.StructField.Name)
		}
	}

	return nil
}
