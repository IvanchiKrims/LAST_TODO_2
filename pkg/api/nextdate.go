package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const isoDate = "20060102"

var (
	ErrBadRepeat      = errors.New("invalid or missing repeat rule")
	ErrInvalidDate    = errors.New("invalid date format, expected YYYYMMDD")
	ErrIntervalTooBig = errors.New("day interval exceeds maximum of 400")
	ErrUnsupported    = errors.New("unsupported repeat rule (only d <N> and y)")
)

func NextDate(now time.Time, dstart, repeat string) (string, error) {
	start, err := time.Parse(isoDate, dstart)
	if err != nil {
		return "", ErrInvalidDate
	}
	if repeat == "" {
		return "", ErrBadRepeat
	}

	var step func(time.Time) (time.Time, error)

	repeat = strings.TrimSpace(repeat)

	if repeat == "y" {
		step = func(t time.Time) (time.Time, error) {
			return t.AddDate(1, 0, 0), nil
		}
	} else if strings.HasPrefix(repeat, "d") {
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", ErrBadRepeat
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", ErrBadRepeat
		}
		if days < 1 || days > 400 {
			return "", ErrIntervalTooBig
		}
		step = func(t time.Time) (time.Time, error) {
			return t.AddDate(0, 0, days), nil
		}
	} else {
		return "", ErrUnsupported
	}

	current := start
	for {
		next, _ := step(current)
		if next.After(now) {
			return next.Format(isoDate), nil
		}
		current = next
	}
}
