package dv

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator/pkg/cv"
)

// General Email Regex (RFC 5322 Official Standard)
// see: http://emailregex.com/
const EmailRegexString = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"

type EmailError string

func EmailErrorf(format string, a ...interface{}) EmailError {
	return EmailError(fmt.Sprintf(format, a...))
}

func (r EmailError) Error() string {
	return string(r)
}

func Email() *cv.CustomValidator {
	emailTagString := "email"
	emailTagRegex := regexp.MustCompile(emailTagString)

	customValidator := cv.NewCustomValidator("email", emailTagRegex, ValidateEmail, cv.NewCustomValidatorConfig().FailForNilValue())
	return customValidator
}

func ValidateEmail(ctx context.Context, f *cv.Field, vCtx *cv.ValidationContext) error {
	value := f.Value
	kind := value.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if f.Value.IsNil() {
			return EmailErrorf("email field %v is nil", f.StructField.Name)
		}

		value = value.Elem()
		kind = value.Kind()
	}

	if kind != reflect.String {
		return EmailErrorf("email field %v cannot be converted to string", f.StructField.Name)
	}

	if value.IsZero() {
		return EmailErrorf("email field %v has zero value", f.StructField.Name)
	}

	emailRegex := regexp.MustCompile(EmailRegexString)
	email := value.String()

	isEmail := emailRegex.MatchString(email)
	if !isEmail {
		return EmailErrorf("email field %v is no valid email", f.StructField.Name)
	}

	return nil
}
