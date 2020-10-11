package dv

import (
	"context"
	"fmt"
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

func ValidateNonZero(ctx context.Context, f *cv.Field) error {
	if f.Value.IsZero() {
		return ZeroErrorf("non-zero field %v has zero value", f.StructField.Name)
	}
	return nil
}
