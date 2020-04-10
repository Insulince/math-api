package operate

import "fmt"

func Multiply(arguments []float64) (product float64, asString string, err error) {
	multiplier := arguments[0]
	multiplicand := arguments[1]

	product = multiplier * multiplicand

	asString = fmt.Sprintf("%v*%v=%v", multiplier, multiplicand, product)

	return product, asString, nil
}
