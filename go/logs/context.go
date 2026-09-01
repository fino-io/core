package logs

import "context"

type contextFieldsKey struct{}

// WithFields returns a child context carrying fields that are safe to include
// in every log entry emitted with that context.
func WithFields(ctx context.Context, fields ...Field) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(fields) == 0 {
		return ctx
	}
	return context.WithValue(ctx, contextFieldsKey{}, mergeFields(fieldsFromContext(ctx), fields))
}

func fieldsFromContext(ctx context.Context) []Field {
	if ctx == nil {
		return nil
	}
	fields, _ := ctx.Value(contextFieldsKey{}).([]Field)
	return fields
}

// mergeFields keeps field order stable while allowing later values to override
// fields with the same key.
func mergeFields(base, extra []Field) []Field {
	fields := append([]Field(nil), base...)
	for _, field := range extra {
		replaced := false
		for i := range fields {
			if fields[i].Key == field.Key {
				fields[i] = field
				replaced = true
				break
			}
		}
		if !replaced {
			fields = append(fields, field)
		}
	}
	return fields
}
