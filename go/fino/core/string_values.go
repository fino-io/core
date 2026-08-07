package core

import "regexp"

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
	for _, e := range x.Vals {
		if e == element {
			return true
		}
	}
	return false
}

func (x *StringValues) Unique() *StringValues {
	tmp := NewStringValues()
	found := make(map[string]struct{})
	for _, v := range x.Vals {
		if _, ok := found[v]; !ok {
			tmp.Append(v)
			found[v] = struct{}{}
		}
	}
	return tmp
}

func (x *StringValues) Matched(expr string) bool {
	if expr == "" {
		return false
	}

	for _, val := range x.Vals {
		if matched, err := regexp.MatchString(val, expr); matched && err == nil {
			return true
		}
	}

	return false
}

func (x *StringValues) Matches(expr string) *StringValues {
	if expr == "" {
		return NewStringValues()
	}

	tmp := NewStringValues()
	for _, val := range x.Vals {
		if matched, err := regexp.MatchString(val, expr); matched && err == nil {
			tmp.Append(val)
		}
	}

	return tmp
}
