// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dateutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// The regular expression which matches ISO 8601 date format pattern (e.g. 2013-02-08).
var iso8601DateFormatPattern = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2})")

// The regular expression which matches a hh:mm time format pattern (e.g 21:13).
var timeFormatPattern = regexp.MustCompile("\\s(\\d{2}):(\\d{2})(:\\d{2})?")

// ParseIso8601Date parses a ISO 8601 date string (e.g. 2013-02-08 21:13)
// and returns the time value it represents.
func ParseIso8601Date(value string, fallback time.Time) (time.Time, error) {

	// Parse the date component (e.g. "2013-02-08")
	// check if the value matches the ISO 8601 Date pattern
	dateComponents := iso8601DateFormatPattern.FindStringSubmatch(value)
	if len(dateComponents) < 4 {
		return fallback, errors.New(fmt.Sprintf("\"%v\" is not a valid ISO 8601 date", value))
	}

	// parse year
	yearString := dateComponents[1]
	yearInt64, parseYearError := strconv.ParseInt(yearString, 10, 16)
	if parseYearError != nil || yearInt64 < 1 || yearInt64 > 9999 {
		return time.Time{}, fmt.Errorf("\"%v\" is not a valid value for a year. Valid values are in the range between 1 and 9999.", yearString)
	}

	// parse month
	monthString := dateComponents[2]
	monthInt64, parseMonthErr := strconv.ParseInt(monthString, 10, 8)
	if parseMonthErr != nil || monthInt64 < 1 || monthInt64 > 12 {
		return time.Time{}, fmt.Errorf("\"%v\" is not a valid value for a month. Valid values are in the range between 1 and 12.", monthString)
	}

	month, parseErrMonth := GetMonth(int(monthInt64))
	if parseErrMonth != nil {
		return time.Time{}, parseErrMonth
	}

	// parse day
	dayString := dateComponents[3]
	dayInt64, parseDayErr := strconv.ParseInt(dayString, 10, 8)
	if parseDayErr != nil || dayInt64 < 1 || dayInt64 > 31 {
		return time.Time{}, fmt.Errorf("\"%v\" is not a valid value for a day. Valid values are in the range between 1 and 31.", dayString)
	}

	// Parse the time component  (e.g. "21:13")
	var (
		hour        int
		minute      int
		second      int
		millisecond int
	)

	// check if the value matches the 24 hour time format pattern

	if timeComponents := timeFormatPattern.FindStringSubmatch(value); len(timeComponents) > 3 {

		// parse hours
		hourString := timeComponents[1]
		hourInt64, parseHourError := strconv.ParseInt(hourString, 10, 16)
		if parseHourError != nil || hourInt64 < 0 || hourInt64 > 23 {
			return time.Time{}, fmt.Errorf("\"%v\" is not a valid value for an hour in a 24h time format. Valid values are in the range between 0 and 23.", hourString)
		}
		hour = int(hourInt64)

		// parse minutes
		minuteString := timeComponents[2]
		minuteInt64, parseMinuteError := strconv.ParseInt(minuteString, 10, 16)
		if parseMinuteError != nil || minuteInt64 < 0 || minuteInt64 > 59 {
			return time.Time{}, fmt.Errorf("\"%v\" is not a valid value for an minute in a 24h time format. Valid values are in the range between 0 and 59.", minuteString)
		}
		minute = int(minuteInt64)

	}

	return time.Date(int(yearInt64), month, int(dayInt64), hour, minute, second, millisecond, time.UTC), nil
}

// GetMonth returns the time.Month value for
// a given integer value in the range between 1 and 12.
func GetMonth(value int) (time.Month, error) {
	switch value {
	case 1:
		return time.January, nil
	case 2:
		return time.February, nil
	case 3:
		return time.March, nil
	case 4:
		return time.April, nil
	case 5:
		return time.May, nil
	case 6:
		return time.June, nil
	case 7:
		return time.July, nil
	case 8:
		return time.August, nil
	case 9:
		return time.September, nil
	case 10:
		return time.October, nil
	case 11:
		return time.November, nil
	case 12:
		return time.December, nil
	}

	return time.January, fmt.Errorf("\"%v\" is not a valid value for a month.", value)
}
