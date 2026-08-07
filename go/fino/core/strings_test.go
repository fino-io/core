package core

import "testing"

func TestQuoteString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{
		{name: "empty", str: "", want: `""`},
		{name: "single quote", str: `"`, want: `"\""`},
		{name: "prefix quote", str: `"foo`, want: `"\"foo"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quote(tt.str); got != tt.want {
				t.Errorf("QuoteString() = |%s|, want |%s|", got, tt.want)
			}
		})
	}
}
