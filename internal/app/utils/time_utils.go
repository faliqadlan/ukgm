package utils

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

// UniversalTimeParser parse date with every possible format
func UniversalTimeParser(source string) (*time.Time, error) {
	date, err := time.Parse(time.ANSIC, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.UnixDate, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RubyDate, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC822, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC822Z, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC850, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC1123, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC1123Z, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC3339, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse(time.RFC3339Nano, source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02T15:04:05Z", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2006-01-02", source)
	if err == nil {
		return &date, nil
	}
	date, err = time.Parse("2-Jan-06", source)
	if err == nil {
		return &date, nil
	}

	return nil, fmt.Errorf("Requested String is not in Time format")
}

func StartOfDay(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return date
}

func EndOfDay(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	return date
}

func StartOfMonth(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), now.New(t).BeginningOfMonth().Day(), 0, 0, 0, 0, t.Location())

	return date
}

func EndOfMonth(t time.Time) time.Time {
	date := time.Date(t.Year(), t.Month(), now.New(t).EndOfMonth().Day(), 23, 59, 59, 59, t.Location())

	return date
}

func StartOfYear(t time.Time) time.Time {
	parse := now.New(t).BeginningOfYear()
	date := time.Date(t.Year(), parse.Month(), parse.Day(), 0, 0, 0, 0, t.Location())

	return date
}

func EndOfYear(t time.Time) time.Time {
	parse := now.New(t).EndOfYear()
	date := time.Date(t.Year(), parse.Month(), parse.Day(), 23, 59, 59, 59, t.Location())

	return date
}

func ChangeTimeZoneToUTC(t *time.Time) (err error) {
	_, offset := t.Zone()
	offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", offset))

	*t = t.Add(offsetDuration).UTC()
	return
}

func ChangeUTCToLocation(t *time.Time, locationString ...string) {
	location := time.Local
	if len(locationString) == 1 {
		temp, err := time.LoadLocation(locationString[0])
		if err == nil {
			location = temp
		}
	}
	_, offset := time.Now().In(location).Zone()
	offset *= -1
	offsetDuration, _ := time.ParseDuration(fmt.Sprintf("%ds", offset))

	*t = t.Add(offsetDuration).In(location)
	return
}

func ChangeTimezoneToLocation(t *time.Time, locationString ...string) {
	ChangeTimeZoneToUTC(t)
	ChangeUTCToLocation(t, locationString...)
}

func DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
