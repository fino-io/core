package core

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
)

// Value Implement driver.Valuer and sql.Scanner interfaces on Duration
func (x *Duration) Value() (driver.Value, error) {
	if x != nil {
		return x.ToSeconds(), nil
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
		if math.IsNaN(duration) || math.IsInf(duration, 0) {
			return fmt.Errorf("Duration.Scan: non-finite seconds")
		}
		x.FromSeconds(duration)
	case []byte:
		return x.Scan(string(duration))
	case string:
		seconds, err := strconv.ParseFloat(duration, 64)
		if err != nil {
			return err
		}
		return x.Scan(seconds)
	default:
		return fmt.Errorf("could not decode type %T -> %T", src, x)
	}
	return nil
}

func (x *Duration) GormDataType() string {
	return "float"
}
