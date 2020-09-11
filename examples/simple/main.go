package main

import (
	"context"
	"fmt"

	"github.com/gogo-gadget/validator"
)

type testStruct struct {
	name   string `validator:"required;gte=10"`
	nested *nestedStruct
}

type nestedStruct struct {
	nestedName   string
	doubleNested *doubleNestedStruct
}

type doubleNestedStruct struct {
	doubleNestedField string `validator:"required"`
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
