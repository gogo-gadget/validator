package validator

import dv "github.com/gogo-gadget/validator/dv"

func (v *Validator) RegisterDefaultCustomValidators() {
	v.RegisterCustomValidator(dv.NonNil())
}
