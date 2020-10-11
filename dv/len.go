package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gogo-gadget/validator/pkg/cv"
)

type LenError string

func LenErrorf(format string, a ...interface{}) LenError {
	return LenError(fmt.Sprintf(format, a...))
}

func (r LenError) Error() string {
	return string(r)
}

func Len() *cv.CustomValidator {
	lenTagString := `len\([1-9][0-9]*\)`
	lenTagRegex := regexp.MustCompile(lenTagString)

	customValidator := cv.NewCustomValidator("len", lenTagRegex, ValidateLen, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

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
