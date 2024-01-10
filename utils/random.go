package utils

import (
	"fmt"
	"math/rand"
	"strconv"
)

func GenerateCode(digits int) int {

	code := ""
	for i := 0; i < digits; i++ {
		digit := rand.Int()

		code += fmt.Sprint(digit)
	}
	num, _ := strconv.Atoi(code)
	return num
}
