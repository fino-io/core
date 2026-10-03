package core

import (
	"regexp"
	"slices"
)

func NewStringValues(vals ...string) *StringValues {
	return &StringValues{Vals: vals}
}

func (x *StringValues) Append(vals ...string) *StringValues {
	if x != nil {
		x.Vals = append(x.Vals, vals...)
	}
	return x
}

func (x *StringValues) ToArray() any {
	if x != nil {
		return x.Vals
	}
	return []string{}
}

func (x *StringValues) Contains(element string) any {
	return slices.Contains(x.GetVals(), element)
}

func (x *StringValues) Unique() *StringValues {
	tmp := NewStringValues()
	found := make(map[string]struct{})
	for _, v := range x.GetVals() {
		if _, ok := found[v]; !ok {
			tmp.Append(v)
			found[v] = struct{}{}
		}
	}
	return tmp
}

func (x *StringValues) Matched(expr string) bool {
	pattern, err := regexp.Compile(expr)
	if expr == "" || err != nil {
		return false
	}

	return slices.ContainsFunc(x.GetVals(), pattern.MatchString)
}

func (x *StringValues) Matches(expr string) *StringValues {
	tmp := NewStringValues()
	pattern, err := regexp.Compile(expr)
	if expr == "" || err != nil {
		return tmp
	}
	for _, val := range x.GetVals() {
		if pattern.MatchString(val) {
			tmp.Append(val)
		}
	}

	return tmp
}
