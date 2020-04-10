package util

import "strconv"

func IsNumber(str string) (isInteger bool) {
	if len(str) == 0 {
		return false
	}

	_, err := strconv.ParseFloat(str, 64)

	return err == nil
}
