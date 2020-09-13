package main

import (
	"context"
	"fmt"
	"log"

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
	doubleNestedField string `validator:"non-nil"`
}

func main() {
	ts := &testStruct{
		name: "top-level",
		nested: &nestedStruct{
			doubleNested: &doubleNestedStruct{
				doubleNestedField: "I am not nil.",
			},
		},
	}

	v := validator.NewValidator()

	ctx := context.Background()

	err := v.Validate(ctx, ts)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("hurray, validation succeeded")
}
