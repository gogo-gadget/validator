package main

import (
	"context"
	"regexp"

	"github.com/gogo-gadget/validator"
	"github.com/gogo-gadget/validator/pkg/cv"
)

func main() {
	v := validator.NewValidator()

	customValidator := exampleValidator()
	v.RegisterCustomValidator(customValidator)
}

func exampleValidator() *cv.CustomValidator {
	exampleString := "example"
	exampleRegexp := regexp.MustCompile(exampleString)

	customValidator := cv.NewCustomValidator("example", exampleRegexp, validateExampleValidator, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

func validateExampleValidator(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	// Validation of the field is placed here
	// ...
	return nil
}
