package magiccfg

import (
	"github.com/go-playground/validator/v10"
)

const (
	minPort = 1023
	maxPort = 65536
)

type ValidationFunction func(validator.FieldLevel) bool

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
