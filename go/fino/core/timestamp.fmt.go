package core

import (
	"errors"
	"time"

	"github.com/araddon/dateparse"
)

func (x *Timestamp) Format() string {
	if x != nil {
		return x.ToTime().UTC().Format(time.RFC3339Nano)
	}
	return ""
}

func (x *Timestamp) ToString() string {
	return x.Format()
}

func ParseTimestamp(value string) (*Timestamp, error) {
	ts := &Timestamp{}
	if err := ts.Parse(value); err != nil {
		return nil, err
	}
	return ts, nil
}

func (x *Timestamp) Parse(value string) error {
	if x == nil {
		return errors.New("timestamp.Parse: nil receiver")
	}
	t, err := dateparse.ParseIn(value, time.UTC)
	if err != nil {
		return err
	}
	parsed := FromTime(t)
	if err := parsed.CheckValid(); err != nil {
		return err
	}
	x.Seconds, x.Nanoseconds = parsed.Seconds, parsed.Nanoseconds
	return nil
}
