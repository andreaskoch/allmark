package date

import (
	"andyk/docs/pattern"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// The regular expression which matches ISO 8601 date format pattern.
var iso8601DateFormatPattern = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2})$")

// ParseIso8601Date parses a ISO 8601 date string (e.g. 2013-02-08 21:13)
// and returns the time value it represents. 
func ParseIso8601Date(value string) (time.Time, error) {

	// check if the value matches the ISO 8601 Date pattern
	isValidIso8601Date, matches := pattern.IsMatch(value, *iso8601DateFormatPattern)

	if !isValidIso8601Date {
		return time.Date(1, 1, 1, 0, 0, 1, 0, time.UTC), errors.New(fmt.Sprintf("\"%v\" is not a valid ISO 8601 date", value))
	}

	// parse year
	yearString := matches[1]
	yearInt64, parseYearError := strconv.ParseInt(yearString, 10, 16)
	if parseYearError != nil {
		log.Panicf("\"%v\" is not a valid year.", yearString)
	}

	// parse month
	monthString := matches[2]
	monthInt64, parseMonthErr := strconv.ParseInt(monthString, 10, 8)
	if parseMonthErr != nil {
		log.Panicf("\"%v\" is not a valid month.", monthString)
	}

	month := GetMonth(int(monthInt64))

	// parse day
	dayString := matches[3]
	dayInt64, parseDayErr := strconv.ParseInt(monthString, 10, 8)
	if parseDayErr != nil {
		log.Panicf("\"%v\" is not a valid day.", dayString)
	}

	return time.Date(int(yearInt64), month, int(dayInt64), 0, 0, 1, 0, time.UTC), nil
}

// GetMonth returns the time.Month value for 
// a given integer value in the range between 1 and 12.
func GetMonth(value int) time.Month {
	switch value {
	case 1:
		return time.January
	case 2:
		return time.February
	case 3:
		return time.March
	case 4:
		return time.April
	case 5:
		return time.May
	case 6:
		return time.June
	case 7:
		return time.July
	case 8:
		return time.August
	case 9:
		return time.September
	case 10:
		return time.October
	case 11:
		return time.November
	case 12:
		return time.December
	}

	panic(fmt.Sprintf("\"%v\" is not a valid value for a month.", value))
}
