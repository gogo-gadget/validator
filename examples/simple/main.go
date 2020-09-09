package main

import (
	"context"
	"fmt"

	"github.com/gogo-gadget/validator"
)

type testStruct struct {
	name string `validator:required,gte=10`
}

func main() {
	ts := testStruct{}

	v := validator.NewValidator()

	ctx := context.Background()

	err := v.Validate(ctx, ts)

	if err != nil{
		fmt.Println("oh no, validation failed")
	}

	fmt.Println("hurray, validation succeeded")
}