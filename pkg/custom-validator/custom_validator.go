package custom_validator

import (
	"context"
	"reflect"
	"regexp"
)

type Field struct {
	StructField reflect.StructField
	Value       reflect.Value
}

type CustomValidationFunc func(ctx context.Context, f Field) error

type CustomValidatorConfig struct {
	// validation will fail if tag is on field of nil ptr
	// or even if tag is nested on some nil ptr
	ShouldFailIfFieldOfNilPtr bool
}

func NewCustomValidatorConfig() *CustomValidatorConfig {
	return &CustomValidatorConfig{}
}

func (cfg *CustomValidatorConfig) WithNilValidation(enabled bool) *CustomValidatorConfig {
	cfg.ShouldFailIfFieldOfNilPtr = enabled
	return cfg
}

type CustomValidator struct {
	ID       string
	TagRegex *regexp.Regexp
	Validate CustomValidationFunc
	Config   *CustomValidatorConfig
}

func NewCustomValidator(id string, tagRegex *regexp.Regexp, validate CustomValidationFunc, cfg *CustomValidatorConfig) *CustomValidator {

	cv := CustomValidator{
		ID:       id,
		TagRegex: tagRegex,
		Validate: validate,
		Config:   cfg,
	}

	return &cv
}
