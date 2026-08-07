package core

import (
	"strconv"
	"strings"
)

const (
	DoubleQuote = `"`
	SingleQuote = `'`
)

func IsQuotedString(str, quote string) bool {
	return strings.HasPrefix(str, quote) && strings.HasSuffix(str, quote)
}

func Quote(str string) string {
	return strconv.Quote(str)
}

func Unquote(str string) (string, error) {
	return strconv.Unquote(str)
}

func QuoteString(str string) string {
	if strings.HasPrefix(str, DoubleQuote) {
		if strings.HasSuffix(str, DoubleQuote) {
			return str
		} else {
			return `"\"` + str[1:] + DoubleQuote
		}
	} else {
		if strings.HasSuffix(str, DoubleQuote) {
			return DoubleQuote + str[:len(str)-1] + `\""`
		} else {
			return DoubleQuote + str + DoubleQuote
		}
	}
}
