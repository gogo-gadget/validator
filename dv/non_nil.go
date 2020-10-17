package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

// Custom Nil Error that will be returned by the Non-Nil Custom Validator
type NilError string

// Creates a new Nil Error by providing a format string and optional parameters
func NilErrorf(format string, a ...interface{}) NilError {
	return NilError(fmt.Sprintf(format, a...))
}

// Returns the error string
// Implements error interface
func (err NilError) Error() string {
	return string(err)
}

// Creates a new non-nil custom validator
func NonNil() *cv.CustomValidator {
	nonNilTagString := "non-nil"
	nonNilTagRegexp := regexp.MustCompile(nonNilTagString)

	customValidator := cv.NewCustomValidator("non-nil", nonNilTagRegexp, ValidateNonNil, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

// Custom validation function for the non-nil custom validator
func ValidateNonNil(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	value := f.Value
	kind := value.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if f.Value.IsNil() {
			return NilErrorf("non-nil field %v is nil", f.StructField.Name)
		}

		value = value.Elem()
		kind = value.Kind()
	}

	return nil
}
