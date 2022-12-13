package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var (
	TIMESTAMP_REGEX = regexp.MustCompile(`^[1-9]\d{3}\d{2}\d{2}T\d{2}\d{2}\d{2}Z$`)
)

// /ptlist?period=1h&tz=Europe/Athens&t1=20210714T204603Z&t2=20210715T123456Z
// isValidPeriod is helper function to parse period and translate it into time.Duration
// Valid periods should be 1h, 1d, 1mo, 1y
func isValidPeriod(p string) bool {
	return p == "1h" || p == "1d" || p == "1mo" || p == "1y"
}

// helper function to check timezone
func parseTimezone(tz string) (string, *ApplicationError) {
	timeZone, err := time.LoadLocation(tz)
	if err != nil {
		return "", &ApplicationError{
			Message:    fmt.Sprintf("could not parse timezone : %s", timeZone),
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return timeZone.String(), nil
}

// CheckInvocationPoint checks if invocation points are in the correct format with the use
// of a regular expression
// ^[1-9]\d{3}\d{2}\d{2}T\d{2}\d{2}\d{2}Z$
func CheckInvocationPoints(t1, t2 string) bool {
	return TIMESTAMP_REGEX.MatchString(t1) &&
		TIMESTAMP_REGEX.MatchString(t2)
}

// CheckInvocationSequence checks if invocation points are in the correct time sequence
func CheckInvocationSequence(t1, t2, layout string) bool {
	ts1, err := time.Parse(layout, t1)
	if err != nil {
		fmt.Println(err)
	}

	ts2, err := time.Parse(layout, t2)
	if err != nil {
		fmt.Println(err)
	}

	return ts1.Before(ts2)
}

// ParseStringToTime receives a string and parses it to time.Time
func ParseStringToTime(layout, invocationPoint string) (*time.Time, *ApplicationError) {
	t1, err := time.Parse(layout, invocationPoint)
	if err != nil {
		return nil, &ApplicationError{
			Message:    "cannot parse invocation points",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return &t1, nil
}

const UTC_FORM = "20060102T150405Z"

// parseInvocationPoints checks invocation points and calculates the timestamps, if any
func parseInvocationPoints(t1, t2 string, period string) ([]string, *ApplicationError) {

	if CheckInvocationPoints(t1, t2) && CheckInvocationSequence(t1, t2, UTC_FORM) {
		return calculateTimestamps(t1, t2, period)
	} else {
		return nil, &ApplicationError{
			Message:    "invocation points do not follow the correct format",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}
}

// calculateTimestampsPerHour appends the timestamps into the slice
func calculateTimestamps(t1, t2 string, period string) ([]string, *ApplicationError) {
	ip1, err := ParseStringToTime(UTC_FORM, t1)
	if err != nil {
		return nil, err
	}

	ip2, err := ParseStringToTime(UTC_FORM, t2)
	if err != nil {
		return nil, err
	}

	return GetTimestamps(ip1, ip2, period)
}

func GetTimestamps(ip1, ip2 *time.Time, period string) ([]string, *ApplicationError) {
	timestamps := []string{}

	switch period {
	case "1h":
		timestamps = getHourlyTimestamps(ip1, ip2, timestamps)
	case "1d":
		timestamps = getDailyTimestamps(ip1, ip2, timestamps)
	case "1mo":
		timestamps = getMonthlyTimestamps(ip1, ip2, timestamps)
	case "1y":
		timestamps = getAnnuallyTimestamps(ip1, ip2, timestamps)
	default:
		return nil, &ApplicationError{
			Message:    "could not parse period",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return timestamps, nil
}

func getAnnuallyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(1, 0, 0)
	}
	return timestamps
}

func getMonthlyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, 1, 0)
	}
	return timestamps
}

func getDailyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, 0, 1)
	}
	return timestamps
}

func getHourlyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.Add(time.Hour)
	}
	return timestamps
}
