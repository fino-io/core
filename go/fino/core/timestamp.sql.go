package core

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Value implements driver.Valuer for Timestamp.
func (x *Timestamp) Value() (driver.Value, error) {
	if x != nil {
		if err := x.CheckValid(); err != nil {
			return nil, err
		}
		return x.ToTime(), nil
	}

	return nil, nil
}

func (x *Timestamp) Scan(src any) error {
	if x == nil {
		return fmt.Errorf("timestamp.Scan: nil receiver")
	}

	switch bs := src.(type) {
	case nil:
		x.Seconds, x.Nanoseconds = 0, 0
	case []byte:
		return x.Parse(string(bs))
	case string:
		return x.Parse(bs)
	case time.Time:
		parsed := FromTime(bs)
		if err := parsed.CheckValid(); err != nil {
			return err
		}
		x.Seconds, x.Nanoseconds = parsed.Seconds, parsed.Nanoseconds
	default:
		return fmt.Errorf("could not decode type %T -> %T", src, x)
	}

	return nil
}

func (x *Timestamp) GormDataType() string {
	return "time"
}
