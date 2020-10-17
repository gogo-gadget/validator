package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gogo-gadget/validator/pkg/cv"
)

// LenError is a custom error that will be returned by the len custom validator
type LenError string

// LenErrorf creates a new len error by providing a format string and optional parameters
func LenErrorf(format string, a ...interface{}) LenError {
	return LenError(fmt.Sprintf(format, a...))
}

// Error returns the error message string
// Implements error interface
func (err LenError) Error() string {
	return string(err)
}

// Len creates a new len custom validator
func Len() *cv.CustomValidator {
	lenTagString := `len\([1-9][0-9]*\)`
	lenTagRegex := regexp.MustCompile(lenTagString)

	customValidator := cv.NewCustomValidator("len", lenTagRegex, ValidateLen, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

// ValidateLen is a custom validation function for the len custom validator
func ValidateLen(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	value := f.Value
	kind := value.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if f.Value.IsNil() {
			return LenErrorf("len field %v is nil", f.StructField.Name)
		}

		value = value.Elem()
		kind = value.Kind()
	}

	var length int
	switch kind {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.String:
		length = value.Len()
	default:
		return LenErrorf("len field %v is of kind %v", f.StructField.Name, kind.String())
	}

	tagLength, err := strconv.Atoi(vCtx.SubTag[4 : len(vCtx.SubTag)-1])
	if err != nil {
		return err
	}

	if length != tagLength {
		return LenErrorf("len field %v has length %v, but should have length %v", f.StructField.Name, length, tagLength)
	}

	return nil
}
