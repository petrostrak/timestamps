package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	TIMESTAMP_REGEX = regexp.MustCompile(`^[1-9]\d{3}\d{2}\d{2}T\d{2}\d{2}\d{2}Z$`)
)

// isValidPeriod checks whether the period query param is valid.
// Valid periods should be 1h, 1d, 1mo, 1y
func isValidPeriod(p string) bool {
	return p == "1h" || p == "1d" || p == "1mo" || p == "1y"
}

// parseTimezone checks whether the timezone is valid.
func parseTimezone(tz string) (string, *ApplicationError) {
	timeZone, err := time.LoadLocation(tz)
	if err != nil {
		return "", &ApplicationError{
			Message:    "could not parse timezone",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return timeZone.String(), nil
}

// CheckInvocationPoint checks if invocation points are in the correct format.
func checkInvocationPoints(t1, t2 string) bool {
	return TIMESTAMP_REGEX.MatchString(t1) &&
		TIMESTAMP_REGEX.MatchString(t2)
}

// checkInvocationSequence checks if invocation points are in the correct sequence.
func checkInvocationSequence(t1, t2, layout string) bool {
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

// parseStringToTime receives a string and parses it to time.Time
func parseStringToTime(layout, invocationPoint string) (*time.Time, *ApplicationError) {
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

// The layout of the timestamp
const UTC_FORM = "20060102T150405Z"

// parseInvocationPoints checks invocation points and calculates the timestamps, if any.
func parseInvocationPoints(t1, t2 string, period string) ([]string, *ApplicationError) {

	if checkInvocationPoints(t1, t2) && checkInvocationSequence(t1, t2, UTC_FORM) {
		ip1, err := parseStringToTime(UTC_FORM, t1)
		if err != nil {
			return nil, err
		}

		ip2, err := parseStringToTime(UTC_FORM, t2)
		if err != nil {
			return nil, err
		}
		return getTimestamps(ip1, ip2, period)
	} else {
		return nil, &ApplicationError{
			Message:    "cannot parse invocation points",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}
}

func getTimestamps(ip1, ip2 *time.Time, period string) ([]string, *ApplicationError) {
	timestamps := []string{}

	n, err := strconv.Atoi(string(period[0]))
	if err != nil {
		return nil, &ApplicationError{
			Message:    "could not parse period",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	switch string(period[1:]) {
	case "h":
		timestamps = getHourlyTimestamps(ip1, ip2, timestamps, n)
	case "d":
		timestamps = getDailyTimestamps(ip1, ip2, timestamps, n)
	case "mo":
		timestamps = getMonthlyTimestamps(ip1, ip2, timestamps, n)
	case "y":
		timestamps = getAnnuallyTimestamps(ip1, ip2, timestamps, n)
	default:
		return nil, &ApplicationError{
			Message:    "could not parse period",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return timestamps, nil
}

func getAnnuallyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string, n int) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(n, 0, 0)
	}
	return timestamps
}

func getMonthlyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string, n int) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, n, 0)
	}
	return timestamps
}

func getDailyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string, n int) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, 0, n)
	}
	return timestamps
}

func getHourlyTimestamps(ip1 *time.Time, ip2 *time.Time, timestamps []string, n int) []string {
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.Add(time.Duration(n) * time.Hour)
	}
	return timestamps
}
