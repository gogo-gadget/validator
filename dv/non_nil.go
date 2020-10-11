package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type NilError string

func NilErrorf(format string, a ...interface{}) NilError {
	return NilError(fmt.Sprintf(format, a...))
}

func (r NilError) Error() string {
	return string(r)
}

func NonNil() *cv.CustomValidator {
	nonNilTagString := "non-nil"
	nonNilTagRegexp := regexp.MustCompile(nonNilTagString)

	customValidator := cv.NewCustomValidator("non-nil", nonNilTagRegexp, ValidateNonNil, cv.NewCustomValidatorConfig().WithNilValidation(true))
	return customValidator
}

func ValidateNonNil(ctx context.Context, f *cv.Field) error {
	kind := f.Value.Kind()

	switch kind {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// the fields value can only be nil if it is an interface, pointer, map, slice, chan, func or unsafe pointer
		if f.Value.IsNil() {
			return NilErrorf("NonNil field %v is nil", f.StructField.Name)
		}
	}

	return nil
}
