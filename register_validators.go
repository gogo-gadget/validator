package validator

import "github.com/gogo-gadget/validator/dv"

// Registers the default custom validators on the validator instance
func (v *Validator) RegisterDefaultCustomValidators() {
	v.RegisterCustomValidator(dv.NonNil())
	v.RegisterCustomValidator(dv.NonZero())
	v.RegisterCustomValidator(dv.Required())
	v.RegisterCustomValidator(dv.Email())
	v.RegisterCustomValidator(dv.Len())
}
