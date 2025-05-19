package magiccfg

import (
	"github.com/go-playground/validator/v10"
)

const (
	minPort = 1023
	maxPort = 65536
)

// Validation function is a type that the magiccfg validation system can call
// to provide custom validation on a field.  Because magiccfg uses
// github.com/go-playground/validator under the hood, more details regarding
// custom validation functions can be found in the documentation to that package.
type ValidationFunction func(validator.FieldLevel) bool

// PortValidator is a validation function that ensures a configuration field
// is populated with a valid, non-system port (i.e. doesn't use 1-1023 and
// doesn't exceed 65535).
func PortValidator(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	if value > minPort && value < maxPort {
		return true
	}
	return false
}
