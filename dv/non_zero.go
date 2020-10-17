package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

// ZeroError is a custom error that will be returned by the non-zero custom validator
type ZeroError string

// ZeroErrorf creates a new zero error by providing a format string and optional parameters
func ZeroErrorf(format string, a ...interface{}) ZeroError {
	return ZeroError(fmt.Sprintf(format, a...))
}

// Error returns the error message string
// Implements error interface
func (err ZeroError) Error() string {
	return string(err)
}

// NonZero creates a new non-zero custom validator
func NonZero() *cv.CustomValidator {
	nonZeroTagString := "non-zero"
	nonZeroTagRegexp := regexp.MustCompile(nonZeroTagString)

	customValidator := cv.NewCustomValidator("non-zero", nonZeroTagRegexp, ValidateNonZero, cv.NewCustomValidatorConfig())
	return customValidator
}

// ValidateNonZero is a custom validation function for the non-zero custom validator
func ValidateNonZero(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	value := f.Value
	kind := value.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if f.Value.IsNil() {
			return ZeroErrorf("non-zero field %v is nil", f.StructField.Name)
		}

		value = value.Elem()
		kind = value.Kind()
	}

	if f.Value.IsZero() {
		return ZeroErrorf("non-zero field %v has zero value", f.StructField.Name)
	}
	return nil
}
