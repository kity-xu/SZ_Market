package controllers

import "strconv"

func PackAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func PackItoa(i int) string {
	return strconv.Itoa(i)
}
