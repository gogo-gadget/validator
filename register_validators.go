package validator

import "github.com/gogo-gadget/validator/dv"

func (v *Validator) RegisterDefaultCustomValidators() {
	v.RegisterCustomValidator(dv.NonNil())
	v.RegisterCustomValidator(dv.NonZero())
	v.RegisterCustomValidator(dv.Required())
	v.RegisterCustomValidator(dv.Email())
	v.RegisterCustomValidator(dv.Len())
}
