# Gogo-Gadget Validator

Gogo-Gadget Validator is a simple struct validator based on field tags with the following features:

- Struct validation by setting validator tags on fields
- Customizable nil pointer validation
- Custom validations can be easily registered
- Context forwarding to be able to add context based validation
- Logical Operators `&&`, `||` and `!` for tags
- Conditional Expressions `if(...)then(...) elif(...)then(...) else(...)` for tags

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

Have a look at the [example](/examples/simple/main.go) below:
```go
package main

import (
	"context"
	"log"

	"github.com/gogo-gadget/validator"
)

type testStruct struct {
	name string `validator:"required && len(9)"`
}

func main() {
	ts := &testStruct{
		name: "top-level",
	}

	v := validator.NewValidator()

	ctx := context.Background()

	err := v.Validate(ctx, ts)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("hurray, validation succeeded")
}
```

## Tag Syntax
The validator tag syntax contains rules for logical operators and conditional expressions. This implies that certain
combinations of characters should not be used in custom validation tag regular expressions to guarantee the correct behavior of the validation.

### Rules
Custom validation tags:  

- should **not** start with: `if(`, `!`
- should **always** include the same number of opening `(` and closing `)` brackets.
- should **not** include any whitespace.
 
### Logical Operators and Conditional Expressions
Negate a validation by placing a `!` in front of the validation
```go
type testStruct struct {
	name   string `validator:"!email"`
}
```

Chain validations by the `&&` or `||` logical operators
```go
type testStruct struct {
	name   string `validator:"required && len(7)"`
}
```

Put a validation into brackets `(...)` to define order of operations
```go
type testStruct struct {
	email   string `validator:"len(28) && (non-nil || email)"`
}
```

Use conditional expressions by writing tags of the form `if(...)then(...) elif(...)then(...) else(...)`
```go
type testStruct struct {
	email   string `validator:"if(email)then(len(28)) elif(required)then(len(10)) else(non-nil)"`
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

Have a look at the [example](/examples/custom-validator/main.go) below:
```go
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
```

Since the id will be used for the registration it allows a regular expression for the field tag to be used multiple times.
That does also imply that if one registers two custom validators with the same id, only the last registered will be used.

## Contribution
Feel free to contribute and e.g. add useful custom validators by opening pull requests.