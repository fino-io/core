package core

import (
	"fmt"
	"time"
)

func (x *Date) CheckValid() error {
	if x == nil {
		return fmt.Errorf("date: nil value")
	}
	if x.Year < 1 || x.Year > 9999 || x.Month < 1 || x.Month > 12 || x.Day < 1 || x.Day > 31 {
		return fmt.Errorf("invalid date: %d-%d-%d", x.Year, x.Month, x.Day)
	}
	t := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, time.UTC)
	if t.Day() != int(x.Day) {
		return fmt.Errorf("invalid date: %d-%d-%d", x.Year, x.Month, x.Day)
	}
	return nil
}

func (x *Date) FromTime(t time.Time) *Date {
	if x != nil {
		x.Year, x.Month, x.Day = int64ToInt32(int64(t.Year())), int32(t.Month()), int32(t.Day())
	}
	return x
}

func (x *Date) ToTime(location *time.Location) (time.Time, error) {
	if err := x.CheckValid(); err != nil {
		return time.Time{}, err
	}
	if location == nil {
		location = time.UTC
	}
	t := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, location)
	if t.Year() != int(x.Year) || t.Month() != time.Month(x.Month) || t.Day() != int(x.Day) {
		return time.Time{}, fmt.Errorf("date: nonexistent local date in %s", location)
	}
	return t, nil
}

func (x *Date) Format() string {
	if x == nil {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", x.Year, x.Month, x.Day)
}

func (x *Date) Parse(value string) error {
	if x == nil {
		return fmt.Errorf("date.Parse: nil receiver")
	}
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}
	parsed := (&Date{}).FromTime(t)
	if err := parsed.CheckValid(); err != nil {
		return err
	}
	x.Year, x.Month, x.Day = parsed.Year, parsed.Month, parsed.Day
	return nil
}

func ParseDate(value string) (*Date, error) {
	x := &Date{}
	if err := x.Parse(value); err != nil {
		return nil, err
	}
	return x, nil
}

func (x *TimeOfDay) CheckValid() error {
	if x == nil {
		return fmt.Errorf("timeOfDay: nil value")
	}
	if x.Hours < 0 || x.Hours > 23 || x.Minutes < 0 || x.Minutes > 59 || x.Seconds < 0 || x.Seconds > 59 || x.Nanoseconds < 0 || x.Nanoseconds >= 1000000000 {
		return fmt.Errorf("invalid time of day: %d:%d:%d.%d", x.Hours, x.Minutes, x.Seconds, x.Nanoseconds)
	}
	return nil
}

func (x *TimeOfDay) FromTime(t time.Time) *TimeOfDay {
	if x != nil {
		x.Hours, x.Minutes, x.Seconds = int32(t.Hour()), int32(t.Minute()), int32(t.Second())
		x.Nanoseconds = int32(t.Nanosecond())
	}
	return x
}

func (x *TimeOfDay) ToDuration() (time.Duration, error) {
	if err := x.CheckValid(); err != nil {
		return 0, err
	}
	return time.Duration(x.Hours)*time.Hour + time.Duration(x.Minutes)*time.Minute + time.Duration(x.Seconds)*time.Second + time.Duration(x.Nanoseconds), nil
}

func (x *TimeOfDay) Format() string {
	if x == nil {
		return ""
	}
	return time.Date(1, 1, 1, int(x.Hours), int(x.Minutes), int(x.Seconds), int(x.Nanoseconds), time.UTC).Format("15:04:05.999999999")
}

func (x *TimeOfDay) Parse(value string) error {
	if x == nil {
		return fmt.Errorf("timeOfDay.Parse: nil receiver")
	}
	t, err := time.Parse("15:04:05.999999999", value)
	if err != nil {
		return err
	}
	x.FromTime(t)
	return nil
}

func ParseTimeOfDay(value string) (*TimeOfDay, error) {
	x := &TimeOfDay{}
	if err := x.Parse(value); err != nil {
		return nil, err
	}
	return x, nil
}

// Location resolves an IANA name first. Without a name, Offset is seconds east
// of UTC. Nil means UTC. Fixed offsets must be whole minutes within ±24 hours.
func (x *TimeZone) Location() (*time.Location, error) {
	if x == nil {
		return time.UTC, nil
	}
	if x.Name != "" {
		return time.LoadLocation(x.Name)
	}
	if x.Offset <= -86400 || x.Offset >= 86400 || x.Offset%60 != 0 {
		return nil, fmt.Errorf("invalid timezone offset: %d", x.Offset)
	}
	return time.FixedZone("", int(x.Offset)), nil
}

func (x *TimeZone) CheckValid() error {
	_, err := x.Location()
	return err
}

func (x *DateTime) FromTime(t time.Time) *DateTime {
	if x != nil {
		x.Year, x.Month, x.Day = int64ToInt32(int64(t.Year())), int32(t.Month()), int32(t.Day())
		x.Hour, x.Minute, x.Seconds = int32(t.Hour()), int32(t.Minute()), int32(t.Second())
		x.Nanoseconds = int32(t.Nanosecond())
		_, offset := t.Zone()
		name := t.Location().String()
		location, err := time.LoadLocation(name)
		// A fixed zone can share an IANA name, and DST folds can have two
		// instants for one wall time. Keep the offset if the name loses fidelity.
		if err != nil || !time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location).Equal(t) {
			name = ""
		}
		x.TimeZone = &TimeZone{Name: name, Offset: int64ToInt32(int64(offset))}
	}
	return x
}

func (x *DateTime) ToTime() (time.Time, error) {
	if x == nil {
		return time.Time{}, fmt.Errorf("dateTime: nil value")
	}
	if err := (&Date{Year: x.Year, Month: x.Month, Day: x.Day}).CheckValid(); err != nil {
		return time.Time{}, err
	}
	if err := (&TimeOfDay{Hours: x.Hour, Minutes: x.Minute, Seconds: x.Seconds, Nanoseconds: x.Nanoseconds}).CheckValid(); err != nil {
		return time.Time{}, err
	}
	location, err := x.TimeZone.Location()
	if err != nil {
		return time.Time{}, err
	}
	t := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), int(x.Hour), int(x.Minute), int(x.Seconds), int(x.Nanoseconds), location)
	// time.Date normalizes nonexistent local times during a DST transition.
	if t.Year() != int(x.Year) || t.Month() != time.Month(x.Month) || t.Day() != int(x.Day) || t.Hour() != int(x.Hour) || t.Minute() != int(x.Minute) || t.Second() != int(x.Seconds) || t.Nanosecond() != int(x.Nanoseconds) {
		return time.Time{}, fmt.Errorf("dateTime: nonexistent local time in %s", location)
	}
	return t, nil
}

func (x *DateTime) CheckValid() error {
	_, err := x.ToTime()
	return err
}

func (x *DateTime) Format() string {
	t, err := x.ToTime()
	if err != nil {
		return ""
	}
	return t.Format(time.RFC3339Nano)
}

func (x *DateTime) Parse(value string) error {
	if x == nil {
		return fmt.Errorf("dateTime.Parse: nil receiver")
	}
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return err
	}
	parsed := (&DateTime{}).FromTime(t)
	if err := parsed.CheckValid(); err != nil {
		return err
	}
	x.Year, x.Month, x.Day = parsed.Year, parsed.Month, parsed.Day
	x.Hour, x.Minute, x.Seconds = parsed.Hour, parsed.Minute, parsed.Seconds
	x.Nanoseconds, x.TimeZone = parsed.Nanoseconds, parsed.TimeZone
	return nil
}

func ParseDateTime(value string) (*DateTime, error) {
	x := &DateTime{}
	if err := x.Parse(value); err != nil {
		return nil, err
	}
	return x, nil
}
