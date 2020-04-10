package operate

import "fmt"

func Subtract(arguments []float64) (difference float64, asString string, err error) {
	minuend := arguments[0]
	subtrahend := arguments[1]

	difference = minuend - subtrahend

	asString = fmt.Sprintf("%v-%v=%v", minuend, subtrahend, difference)

	return difference, asString, nil
}
