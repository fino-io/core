package core

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	DoubleQuote = `"`
	SingleQuote = `'`
)

func IsQuotedString(str, quote string) bool {
	return quote != "" && len(str) >= 2*len(quote) && strings.HasPrefix(str, quote) && strings.HasSuffix(str, quote)
}

func Quote(str string) string {
	return strconv.Quote(str)
}

func Unquote(str string) (string, error) {
	return strconv.Unquote(str)
}

func QuoteString(str string) string {
	if IsQuotedString(str, DoubleQuote) && json.Valid([]byte(str)) {
		return str
	}
	quoted, _ := json.Marshal(str)
	return string(quoted)
}
