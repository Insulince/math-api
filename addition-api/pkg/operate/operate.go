package operate

import "fmt"

func Add(arguments []float64) (sum float64, asString string, err error) {
	addend1 := arguments[0]
	addend2 := arguments[1]

	sum = addend1 + addend2

	asString = fmt.Sprintf("%v+%v=%v", addend1, addend2, sum)

	return sum, asString, nil
}
