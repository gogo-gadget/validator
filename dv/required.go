package dv

import (
	"context"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

func Required() *cv.CustomValidator {
	requiredTagString := "required"
	requiredTagRegexp := regexp.MustCompile(requiredTagString)

	customValidator := cv.NewCustomValidator("required", requiredTagRegexp, ValidateRequired, cv.NewCustomValidatorConfig().WithNilValidation(true))
	return customValidator
}

func ValidateRequired(ctx context.Context, f *cv.Field) error {
	err := ValidateNonNil(ctx, f)
	if err != nil {
		return err
	}
	err = ValidateNonZero(ctx, f)
	if err != nil {
		return err
	}

	return nil
}
