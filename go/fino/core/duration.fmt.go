package core

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

func (x *Duration) Format() string {
	if x != nil {
		seconds := new(big.Rat).SetFrac(x.totalNanoseconds(), big.NewInt(int64(time.Second)))
		return strings.TrimRight(strings.TrimRight(seconds.FloatString(9), "0"), ".") + "s"
	}
	return ""
}

func (x *Duration) ToString() string {
	return x.Format()
}

func ParseDuration(value string) (*Duration, error) {
	duration := &Duration{}
	if err := duration.Parse(value); err != nil {
		return nil, err
	}
	return duration, nil
}

func (x *Duration) Parse(value string) error {
	if strings.HasSuffix(value, "s") {
		if err := x.parseSeconds(strings.TrimSuffix(value, "s")); err == nil {
			return nil
		}
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	if x == nil {
		return fmt.Errorf("duration.Parse: nil receiver")
	}
	x.FromDuration(d)
	return nil
}
