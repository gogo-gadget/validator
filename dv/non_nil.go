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

	customValidator := cv.NewCustomValidator("non-nil", nonNilTagRegexp, ValidateNonNil, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

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
