package operate

import (
	"fmt"
	"errors"
)

func Divide(arguments []float64) (quotient float64, asString string, err error) {
	dividend := arguments[0]
	divisor := arguments[1]

	if divisor == 0 {
		return 0, "invalid", errors.New("Divisor cannot be zero!\n")
	}

	quotient = dividend / divisor

	asString = fmt.Sprintf("%v/%v=%v", dividend, divisor, quotient)

	return quotient, asString, nil
}
