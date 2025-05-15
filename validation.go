package magiccfg

import (
	"github.com/go-playground/validator/v10"
	"time"
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

func TimeValidator(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if _, err := time.Parse(time.RFC3339, value); err == nil {
		return true
	}

	value += "Z"

	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return false
	}

	return true
}
