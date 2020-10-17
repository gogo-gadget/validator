package cv

import (
	"context"
	"reflect"
	"regexp"
)

// ValidationContext contains information about the current validation.
// Will be forwarded to Custom Validators.
// Contains the current SubTag that is validated.
type ValidationContext struct {
	SubTag string
}

// Field contains information about the field that is validated.
type Field struct {
	// Parent is either the parent field or nil if the field has no parent.
	Parent      *Field
	StructField reflect.StructField
	Value       reflect.Value
}

// CustomValidationFunc is the type of validation function that needs to be provided in custom validator to be run on struct fields
type CustomValidationFunc func(ctx context.Context, f *Field, validationCtx *ValidationContext) error

// CustomValidatorConfig is used to configure a custom validator
type CustomValidatorConfig struct {
	// Validation will fail if tag is on field of nil ptr
	// or even if tag is nested on some nil ptr
	ShouldFailIfFieldOfNilPtr bool
}

// NewCustomValidatorConfig creates a new custom validator configuration
func NewCustomValidatorConfig() *CustomValidatorConfig {
	return &CustomValidatorConfig{}
}

// FailForNilValue configures the custom validator configuration to fail if the field is located on a nil pointer
func (cfg *CustomValidatorConfig) FailForNilValue() *CustomValidatorConfig {
	cfg.ShouldFailIfFieldOfNilPtr = true
	return cfg
}

// CustomValidator is used to run validations on struct field tags
type CustomValidator struct {
	// ID of the Custom Validator
	// This should be unique otherwise the last registered Custom Validator will replace the previous with the same ID
	ID string
	// A regular expression which decides based on a StructFieldTag if the Custom Validation Func should be executed on a StructField.
	TagRegex *regexp.Regexp
	// The Custom Validation Func that should be executed on a Field
	Validate CustomValidationFunc
	// The configuration for the Custom Validator
	Config *CustomValidatorConfig
}

// NewCustomValidator creates a new Custom Validator
func NewCustomValidator(id string, tagRegex *regexp.Regexp, validate CustomValidationFunc, cfg *CustomValidatorConfig) *CustomValidator {
	cv := CustomValidator{
		ID:       id,
		TagRegex: tagRegex,
		Validate: validate,
		Config:   cfg,
	}

	return &cv
}
