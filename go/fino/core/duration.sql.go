package core

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Value Implement driver.Valuer and sql.Scanner interfaces on Duration
func (x *Duration) Value() (driver.Value, error) {
	if x != nil {
		return strings.TrimSuffix(x.Format(), "s"), nil
	}
	return nil, nil
}

func (x *Duration) Scan(src any) error {
	if x == nil {
		return fmt.Errorf("Duration.Scan: nil receiver")
	}

	switch duration := src.(type) {
	case nil:
		x.Seconds, x.Nanoseconds = 0, 0
	case int64:
		x.Seconds, x.Nanoseconds = duration, 0
	case float64:
		return x.FromSeconds(duration)
	case []byte:
		return x.Scan(string(duration))
	case string:
		return x.parseSeconds(duration)
	default:
		return fmt.Errorf("could not decode type %T -> %T", src, x)
	}
	return nil
}

func (x *Duration) GormDataType() string {
	return "text"
}
