package utila

import (
	"strconv"
	"strings"
)

func CapitalPeriod(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToUpper(s[:1]) + s[1:]
	switch s[len(s)-1:] {
	case ".", "!", "?":
		return s
	default:
		return s + "."
	}
}

func OrdinalNumber(i int) string {
	n := strconv.Itoa(i)
	if i >= 11 && i <= 13 {
		return n + "th"
	}
	switch i % 10 {
	case 1:
		return n + "st"
	case 2:
		return n + "nd"
	case 3:
		return n + "rd"
	default:
		return n + "th"
	}
}
