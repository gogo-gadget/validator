package dv

import (
	"context"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

func Required() *cv.CustomValidator {
	requiredTagString := "required"
	requiredTagRegexp := regexp.MustCompile(requiredTagString)

	customValidator := cv.NewCustomValidator("required", requiredTagRegexp, ValidateRequired, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

func ValidateRequired(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	err := ValidateNonNil(ctx, f, vCtx)
	if err != nil {
		return err
	}
	err = ValidateNonZero(ctx, f, vCtx)
	if err != nil {
		return err
	}

	return nil
}
