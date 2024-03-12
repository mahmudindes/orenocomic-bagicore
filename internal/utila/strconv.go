package utila

import "strconv"

func Atou(s string) (uint, error) {
	result, err := strconv.ParseUint(s, 10, 0)
	return uint(result), err
}

func Utoa(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}
