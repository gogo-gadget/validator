package validator

import "github.com/gogo-gadget/validator/dv"

func (v *Validator) RegisterDefaultCustomValidators() {
	v.RegisterCustomValidator(dv.NonNil())
}
