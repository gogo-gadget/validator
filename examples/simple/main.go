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
