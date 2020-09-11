# Gogo-Gadget Validator

Gogo-Gadget Validator is a simple struct validator based on field tags with the following features:

- Struct validation by setting validator tags on fields
- Customizable nil pointer validation
- Custom validations can be easily registered
- Context forwarding to be able to add context based validation

## Setup

Download the module with go get:
```
go get github.com/gogo-gadget/validator
```

Import the module by adding:
```go
import "github.com/gogo-gadget/validator"
```

## Usage

In order to use the validation simply create a new instance of a validator by calling `NewValidator()`.
You can then pass any struct, pointer or interface with underlying structs to its `Validate` function.
An error will be returned if the validation failed and nil otherwise.

```go
package main

import (
	"context"
	"fmt"

	"github.com/gogo-gadget/validator"
)

type testStruct struct {
	name   string `validator:"non-nil"`
}

func main() {
	ts := &testStruct{
		name: "top-level",
	}

	v := validator.NewValidator()

	ctx := context.Background()

	err := v.Validate(ctx, ts)

	if err != nil {
		fmt.Println("oh no, validation failed")
		return
	}

	fmt.Println("hurray, validation succeeded")
}
```

## Tag Syntax
Multiple subtags can be provided on a field by separating them with a semicolon.
e.g.

```go
type testStruct struct {
	name   string `validator:"required;non-nil"`
}
```

## Custom Validation

In order to register a custom validator you need an instance of a validator and a custom validator being registered to it.
In order to create a custom validator one can use the provided factory function `cv.NewCustomValidator`.

The function takes an id, regular expression, validation function and a configuration as parameters.

- The id is mainly used for the registration of the custom validator.
- The regular expression is being used to identify if a field should be validated or not.
- The validation function will be run on a field if the regular expression matched a subtag and potentially return an error.
- The configuration allows e.g. to define if the validation should fail if the field is part of a nil pointer to a struct.

```go
package main

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/gogo-gadget/validator"
    cv "github.com/gogo-gadget/validator/pkg/cv"
)

func main() {
	v := validator.NewValidator()

	customValidator := NonNil()
	v.RegisterCustomValidator(customValidator)
}

func NonNil() *cv.CustomValidator {
	nonNilString := "non-nil"
	nonNilRegexp := regexp.MustCompile(nonNilString)

	customValidator := cv.NewCustomValidator("non-nil", nonNilRegexp, ValidateNonNil, cv.NewCustomValidatorConfig().WithNilValidation(true))
	return customValidator
}

func ValidateNonNil(ctx context.Context, f cv.Field) error {
	kind := f.Value.Kind()

	switch kind {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// the fields value can only be nil if it is an interface, pointer, map, slice, chan, func or unsafe pointer
		if f.Value.IsNil() {
			return fmt.Errorf("NonNil field %v is nil", f.StructField.Name)
		}
	}

	return nil
}
```

Since the id will be used for the registration it allows a regular expression for the field tag to be used multiple times.
That does also imply that if one registers two custom validators with the same id, only the last registered will be used.

## Contribution
Feel free to contribute and e.g. add useful custom validators by opening pull requests.