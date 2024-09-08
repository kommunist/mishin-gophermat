package luhn

import (
	"strconv"
)

func Valid(code []byte) (bool, error) {
	l := len(code)

	sum := 0

	for i := 0; i < l; i++ {
		r := l - i

		val, _ := strconv.Atoi(string(code[i])) // сделатб обработку ошибки

		if r%2 == 0 {
			tw := 2 * val
			if tw > 9 {
				sum += (tw - 9)
			} else {
				sum += tw
			}
		} else {
			sum += val
		}
	}

	return sum%10 == 0, nil
}
