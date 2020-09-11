package validator

import dv "github.com/gogo-gadget/validator/default-validators"

func (v *Validator) RegisterDefaultCustomValidators() {
	v.RegisterCustomValidator(dv.NonNil())
}
