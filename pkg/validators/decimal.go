package validators

import (
	"math"

	"github.com/go-playground/validator/v10"
)

// MaxTwoDecimals validates that a float64 value has at most 2 decimal places
func MaxTwoDecimals(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	multiplied := value * 100
	return math.Abs(multiplied-math.Round(multiplied)) < 0.0001
}
