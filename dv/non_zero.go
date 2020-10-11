package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type ZeroError string

func ZeroErrorf(format string, a ...interface{}) ZeroError {
	return ZeroError(fmt.Sprintf(format, a...))
}

func (r ZeroError) Error() string {
	return string(r)
}

func NonZero() *cv.CustomValidator {
	nonZeroTagString := "non-zero"
	nonZeroTagRegexp := regexp.MustCompile(nonZeroTagString)

	customValidator := cv.NewCustomValidator("non-zero", nonZeroTagRegexp, ValidateNonZero, cv.NewCustomValidatorConfig())
	return customValidator
}

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
