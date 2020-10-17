package validator

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/gogo-gadget/validator/pkg/cv"
)

// Custom Tag Syntax Error that will be returned by the validation if the syntax for the validator tag is not correct
type TagSyntaxError struct {
	error  string
	Fields map[string]interface{}
}

// Creates a new Tag Syntax Error by providing a format string and optional parameters
func SyntaxErrorf(format string, a ...interface{}) *TagSyntaxError {
	return &TagSyntaxError{
		error:  fmt.Sprintf(format, a...),
		Fields: map[string]interface{}{},
	}
}

// Returns the error string
// Implements error interface
func (err *TagSyntaxError) Error() string {
	if len(err.Fields) > 0 {
		//TODO this needs a better formatting
		return fmt.Sprintf("%v: %v", err.error, err.Fields)
	}
	return err.error
}

// Adds custom information of the form "key: value" to a Syntax Error
func (err *TagSyntaxError) WithField(key string, value interface{}) *TagSyntaxError {
	err.Fields[key] = value

	return err
}

// Adds multiple information to a Syntax Error each of the form "key: value"
func (err *TagSyntaxError) WithFields(data map[string]interface{}) *TagSyntaxError {
	for key, val := range data {
		err.Fields[key] = val
	}

	return err
}

// A Validator can be used to validate instances of structs or pointers to structs.
// Uses StructFieldTags of form `validator:"..."` to identify validation rules on the field.
// Contains a map of Custom Validators that will be used for the validation.
type Validator struct {
	CustomValidators map[string]*cv.CustomValidator
}

// Creates a new instance of a validator and registers all provided default custom validators for it.
// Usage:
// 		validator := NewValidator()
//		err := validator.Validate(...)
func NewValidator() *Validator {
	v := &Validator{}

	v.RegisterDefaultCustomValidators()

	return v
}

// Registers a custom validator for the validator.
func (v *Validator) RegisterCustomValidator(customValidator *cv.CustomValidator) {
	if v.CustomValidators == nil {
		v.CustomValidators = map[string]*cv.CustomValidator{}
	}

	v.CustomValidators[customValidator.ID] = customValidator
}

// Validates the provided interface{} and forwards the provided context to all custom validators.
// Returns an error if the validation failed or nil otherwise.
func (v *Validator) Validate(ctx context.Context, i interface{}) error {
	iValue := reflect.ValueOf(i)
	iType := iValue.Type()
	kind := iValue.Kind()

	// if the kind of the provided interface is interface or pointer use its underlying element instead
	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if iValue.IsNil() {
			// fail validators that should fail on a nil ptr
			iType = getUnderlyingType(iType)
			kind = iType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(iType, nil)
			if err != nil {
				return err
			}

			return nil
		}

		iValue = iValue.Elem()
		iType = iValue.Type()
		kind = iValue.Kind()
	}

	if kind != reflect.Struct {
		return fmt.Errorf("validation of kind %v is not supported", kind)
	}

	err := v.validateStruct(ctx, iValue, nil)
	if err != nil {
		return err
	}

	return nil
}

// Should only be used on reflect.Values of kind struct
func (v *Validator) validateStruct(ctx context.Context, structValue reflect.Value, parent *cv.Field) error {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		fieldValue := structValue.Field(i)

		field := &cv.Field{
			Parent:      parent,
			StructField: structField,
			Value:       fieldValue,
		}

		err := v.validateField(ctx, field)
		if err != nil {
			return err
		}
	}

	return nil
}

// Is run on every field and sub field of a struct
func (v *Validator) validateField(ctx context.Context, field *cv.Field) error {
	// Validate Field if it contains a subTag matching a regex of any custom validator
	validatorTag := field.StructField.Tag.Get("validator")
	err := v.runFieldValidation(ctx, field, validatorTag)
	if err != nil {
		return err
	}

	fValue := field.Value
	fType := fValue.Type()
	kind := fValue.Kind()

	// if the kind of the Field is interface or pointer use its underlying element instead
	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		if fValue.IsNil() {
			// fail validators that should fail on a nil ptr
			fType = getUnderlyingType(fType)
			kind = fType.Kind()

			if kind != reflect.Struct {
				// if the kind is not struct there is nothing to be validated
				return nil
			}

			err := v.validateStructNilValidations(fType, field)
			if err != nil {
				return err
			}

			return nil
		}

		fValue = fValue.Elem()
		fType = fValue.Type()
		kind = fValue.Kind()
	}

	// If the Field is not of kind struct there is nothing to be validated anymore
	if kind != reflect.Struct {
		return nil
	}

	// If the Field itself is of kind struct validate the nested struct
	err = v.validateStruct(ctx, fValue, field)
	if err != nil {
		return err
	}

	return nil
}

func (v *Validator) validateStructNilValidations(structType reflect.Type, parent *cv.Field) error {
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)

		field := &cv.Field{
			Parent:      parent,
			StructField: structField,
		}

		err := v.validateFieldNilValidations(field)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) validateFieldNilValidations(field *cv.Field) error {
	structField := field.StructField
	validatorTag := structField.Tag.Get("validator")
	subTags := strings.Split(validatorTag, ";")

	for _, customValidator := range v.CustomValidators {
		for _, subTag := range subTags {
			if customValidator.TagRegex.MatchString(subTag) && customValidator.Config.ShouldFailIfFieldOfNilPtr {
				fullFieldName := getFullFieldName(field)
				return fmt.Errorf("validation failed since validator for regex: %v failed on nil value for Field: %v", customValidator.TagRegex.String(), fullFieldName)
			}
		}
	}

	fType := getUnderlyingType(structField.Type)
	kind := fType.Kind()

	// If the Field is not of kind struct there is nothing to be validated anymore
	if kind != reflect.Struct {
		return nil
	}

	// If the Field itself is of kind struct validate the nested struct
	err := v.validateStructNilValidations(fType, field)
	if err != nil {
		return err
	}

	return nil
}

// StructFieldTag validation

// TODO validate validation tag syntax
func (v *Validator) runFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	strippedTag := removeWhiteSpace(tag)
	return v.executeFieldValidation(ctx, field, strippedTag)
}

// TODO add multiple errors to return
func (v *Validator) executeFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	if strings.HasPrefix(tag, "if(") {
		// tag starts with if( statement
		return v.executeIfFieldValidation(ctx, field, tag)
	} else if strings.HasPrefix(tag, "!") {
		// tag starts with ! statement
		return v.executeNotFieldValidation(ctx, field, tag)
	} else if strings.HasPrefix(tag, "(") {
		// tag starts with statement in brackets
		return v.executeBracketFieldValidation(ctx, field, tag)
	}

	// assume tag starts with validation tag
	return v.executeTagFieldValidation(ctx, field, tag)
}

func (v *Validator) executeIfFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	tagLen := len(tag)

	numOpenBraces := 1
	i := 3
	// look forward until all open braces of if condition are closed
	for numOpenBraces > 0 {
		if tag[i] == '(' {
			numOpenBraces++
		} else if tag[i] == ')' {
			numOpenBraces--
		}
		i++
	}

	if i+7 > tagLen || tag[i:i+5] != "then(" {
		// syntax error
		return SyntaxErrorf("if condition must be followed by then statement").
			WithFields(map[string]interface{}{
				"tag":        tag,
				"field-path": getFullFieldName(field),
			})
	}

	numOpenBraces = 1
	j := i + 5
	// look forward until all open braces of then statement are closed
	for numOpenBraces > 0 {
		if tag[j] == '(' {
			numOpenBraces++
		} else if tag[j] == ')' {
			numOpenBraces--
		}
		j++
	}

	// read to end => only one sub validation
	if j == tagLen {
		conditionErr := v.executeFieldValidation(ctx, field, tag[3:i-1])

		if conditionErr == nil {
			return v.executeFieldValidation(ctx, field, tag[i+5:j-1])
		}
		return nil
	}

	if j+2 >= tagLen {
		// at least 3 characters must follow then statement at this point
		return SyntaxErrorf("then statement must be followed by &&, ||, elif or else statement").
			WithFields(map[string]interface{}{
				"tag":        tag,
				"field-path": getFullFieldName(field),
			})
	}

	k := j
	if j+7 < tagLen {
		// read all elif and possible else statements
		if tag[k:k+5] == "elif(" {
			// read all elif statements
			k = j + 5
			numOpenBraces = 1
			for numOpenBraces > 0 {
				if tag[k] == '(' {
					numOpenBraces++
				} else if tag[k] == ')' {
					numOpenBraces--
				}
				k++
			}

			if k+7 >= tagLen || tag[i:i+5] != "then(" {
				// syntax error
				return SyntaxErrorf("elif condition must be followed by then statement").
					WithFields(map[string]interface{}{
						"tag":        tag,
						"field-path": getFullFieldName(field),
					})
			}

			numOpenBraces = 1
			k = k + 5
			// look forward until all open braces of then statement are closed
			for numOpenBraces > 0 {
				if tag[k] == '(' {
					numOpenBraces++
				} else if tag[k] == ')' {
					numOpenBraces--
				}
				k++
			}
		}

		if k+7 < tagLen {
			if tag[k:k+5] == "else(" {
				// read else statement
				k = k + 5
				numOpenBraces = 1
				for numOpenBraces > 0 {
					if tag[k] == '(' {
						numOpenBraces++
					} else if tag[k] == ')' {
						numOpenBraces--
					}
					k++
				}
			}
		}
	}

	if k == tagLen {
		// read to the end of the tag => no && or || statement followed entire if statement
		conditionErr := v.executeFieldValidation(ctx, field, tag[3:i-1])

		if conditionErr == nil {
			// return execution of first then statement
			return v.executeFieldValidation(ctx, field, tag[i+5:j-1])
		}

		if j+7 >= tagLen {
			// at least 7 characters must follow previous if then statement at this point
			return SyntaxErrorf("then statement must be followed by &&, ||, elif or else statement").
				WithFields(map[string]interface{}{
					"tag":        tag,
					"field-path": getFullFieldName(field),
				})
		}

		if tag[j:j+5] == "else(" {
			// return execution of else statement
			return v.executeFieldValidation(ctx, field, tag[j+5:tagLen-1])
		} else if tag[j:j+5] == "elif(" {
			// return execution starting from elif statement
			return v.executeFieldValidation(ctx, field, tag[j+2:])
		}

		return SyntaxErrorf("then statement must be followed by elif or else statement").
			WithFields(map[string]interface{}{
				"tag":        tag,
				"field-path": getFullFieldName(field),
			})
	}

	if k+2 >= tagLen {
		// at least 3 characters must follow then statement at this point
		return SyntaxErrorf("statement must be followed by && or || statement").
			WithFields(map[string]interface{}{
				"tag":        tag,
				"field-path": getFullFieldName(field),
			})
	}

	// validate logical operators "&&" and "||"
	if tag[k:k+2] == "&&" {
		return v.executeAndFieldValidation(ctx, field, tag[:k], tag[k+2:])
	} else if tag[k:k+2] == "||" {
		return v.executeOrFieldValidation(ctx, field, tag[:k], tag[k+2:])
	}

	return SyntaxErrorf("then statement must be followed by && or ||").
		WithFields(map[string]interface{}{
			"tag":        tag,
			"field-path": getFullFieldName(field),
		})
}

func (v *Validator) executeNotFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	tagLen := len(tag)

	if tagLen < 3 {
		// only one sub validation
		if v.executeFieldValidation(ctx, field, tag[1:]) != nil {
			return nil
		}

		return fmt.Errorf("validation of %v of Field %v failed", tag, getFullFieldName(field))
	}

	numOpenBraces := 0
	i := 1
	// look forward until all open braces are closed and either && or || follows
	// last two characters may not be && or ||
	for (i+2 < tagLen && tag[i:i+2] != "&&" && tag[i:i+2] != "||") || numOpenBraces > 0 {
		if tag[i] == '(' {
			numOpenBraces++
		} else if tag[i] == ')' {
			numOpenBraces--
		}
		i++
	}

	// read to end => only one sub validation
	if i+2 >= tagLen {
		if v.executeFieldValidation(ctx, field, tag[1:]) != nil {
			return nil
		}

		return fmt.Errorf("validation of %v of Field %v failed", tag, getFullFieldName(field))
	}

	// validate logical operators && and ||
	if tag[i:i+2] == "&&" {
		return v.executeAndFieldValidation(ctx, field, tag[:i], tag[i+2:])
	} else if tag[i:i+2] == "||" {
		return v.executeOrFieldValidation(ctx, field, tag[:i], tag[i+2:])
	}

	// should not happen
	return SyntaxErrorf("negation operator must be followed by && or || logical operator").WithFields(map[string]interface{}{
		"tag":        tag,
		"field-path": field,
	})
}

func (v *Validator) executeBracketFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	numOpenBraces := 1
	i := 1

	// look forward until all open braces are closed
	for numOpenBraces > 0 {
		if tag[i] == '(' {
			numOpenBraces++
		} else if tag[i] == ')' {
			numOpenBraces--
		}
		i++
	}

	// read to end => only one sub validation
	if i == len(tag) {
		return v.executeFieldValidation(ctx, field, tag[1:i-1])
	}

	// validate logical operators && and ||
	if tag[i:i+2] == "&&" {
		return v.executeAndFieldValidation(ctx, field, tag[1:i-1], tag[i+2:])
	} else if tag[i:i+2] == "||" {
		return v.executeOrFieldValidation(ctx, field, tag[1:i-1], tag[i+2:])
	}

	// should not happen
	return SyntaxErrorf("closing bracket must be followed by && or || logical operator").
		WithFields(map[string]interface{}{
			"tag":        tag,
			"field-path": getFullFieldName(field),
		})
}

func (v *Validator) executeTagFieldValidation(ctx context.Context, field *cv.Field, tag string) error {
	tagLen := len(tag)

	if tagLen < 3 {
		// only one sub validation
		// the tag has to be a validation tag
		for _, customValidator := range v.CustomValidators {
			if customValidator.TagRegex.MatchString(tag) {
				validationCtx := &cv.ValidationContext{
					SubTag: tag,
				}
				err := customValidator.Validate(ctx, field, validationCtx)
				if err != nil {
					return err
				}
			}
		}

		return fmt.Errorf("validation of %v of Field %v failed", tag, getFullFieldName(field))
	}

	numOpenBraces := 0
	i := 1
	// look forward until all open braces are closed and either && or || follows
	// !if()then() is not supported
	for (i+2 < tagLen && tag[i:i+2] != "&&" && tag[i:i+2] != "||") || numOpenBraces > 0 {
		if tag[i] == '(' {
			numOpenBraces++
		} else if tag[i] == ')' {
			numOpenBraces--
		}
		i++
	}

	// read to end => only one sub validation
	if i+2 >= tagLen {
		// the tag has to be a single validation tag
		for _, customValidator := range v.CustomValidators {
			if customValidator.TagRegex.MatchString(tag) {
				validationCtx := &cv.ValidationContext{
					SubTag: tag,
				}
				err := customValidator.Validate(ctx, field, validationCtx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// validate logical operators && and ||
	if tag[i:i+2] == "&&" {
		return v.executeAndFieldValidation(ctx, field, tag[:i], tag[i+2:])
	} else if tag[i:i+2] == "||" {
		return v.executeOrFieldValidation(ctx, field, tag[:i], tag[i+2:])
	}

	// should not happen
	return SyntaxErrorf("validation tag must be followed by && or || logical operator").
		WithFields(map[string]interface{}{
			"tag":        tag,
			"field-path": getFullFieldName(field),
		})
}

func (v *Validator) executeAndFieldValidation(ctx context.Context, field *cv.Field, tag1, tag2 string) error {
	error1 := v.executeFieldValidation(ctx, field, tag1)
	error2 := v.executeFieldValidation(ctx, field, tag2)

	if error1 != nil || error2 != nil {
		return fmt.Errorf("&& validation of %v and %v of Field %v failed", tag1, tag2, getFullFieldName(field))
	}
	return nil
}

func (v *Validator) executeOrFieldValidation(ctx context.Context, field *cv.Field, tag1, tag2 string) error {
	error1 := v.executeFieldValidation(ctx, field, tag1)
	error2 := v.executeFieldValidation(ctx, field, tag2)

	if error1 != nil && error2 != nil {
		return fmt.Errorf("|| validation of %v or %v of Field %v failed", tag1, tag2, getFullFieldName(field))
	}
	return nil
}

// Utility Methods

// removes the whitespace from a provided string
func removeWhiteSpace(str string) string {
	strippedStr := ""
	for _, r := range str {
		if !unicode.IsSpace(r) {
			strippedStr = fmt.Sprintf("%v%c", strippedStr, r)
		}
	}

	return strippedStr
}

func getFullFieldName(field *cv.Field) string {
	if field == nil {
		return ""
	}

	fullFieldName := field.StructField.Name
	parent := field.Parent
	for parent != nil {
		fullFieldName = fmt.Sprintf("%v.%v", parent.StructField.Name, fullFieldName)
		parent = parent.Parent
	}

	return fullFieldName
}

func getUnderlyingType(rType reflect.Type) reflect.Type {
	kind := rType.Kind()

	for kind == reflect.Interface || kind == reflect.Ptr || kind == reflect.UnsafePointer {
		rType = rType.Elem()
		kind = rType.Kind()
	}

	return rType
}
