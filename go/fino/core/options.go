package core

import "sort"

type Options map[string]any

func NewOptions(kvs ...any) Options {
	options := make(Options)
	return options.SetValues(kvs...)
}

func (o Options) SetValue(key string, value any) Options {
	if o == nil {
		o = make(Options)
	}
	o[key] = value
	return o
}

// SetValues ignores non-string keys and an incomplete trailing pair.
// Retain the returned map when using a nil Options value.
func (o Options) SetValues(kvs ...any) Options {
	for i := 0; i < len(kvs)-1; i += 2 {
		if key, ok := kvs[i].(string); ok {
			o = o.SetValue(key, kvs[i+1])
		}
	}
	return o
}

func (o Options) Merge(options Options) Options {
	for k, v := range options {
		o = o.SetValue(k, v)
	}
	return o
}

func (o Options) KeyValues() []any {
	keys := make([]string, 0, len(o))
	for key := range o {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	kvs := make([]any, 0, len(o)*2)
	for _, key := range keys {
		kvs = append(kvs, key, o[key])
	}
	return kvs
}
