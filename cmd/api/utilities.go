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
func parseTimezone(tz string) (*time.Location, *ApplicationError) {
	timeZone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, &ApplicationError{
			Message:    "could not parse timezone",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}

	return timeZone, nil
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
func parseStringToTime(invocationPoint string) (*time.Time, *ApplicationError) {
	t1, err := time.Parse(UTC_FORM, invocationPoint)
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
func parseInvocationPoints(t1, t2, tz string, period string) ([]string, *ApplicationError) {

	if checkInvocationPoints(t1, t2) && checkInvocationSequence(t1, t2, UTC_FORM) {
		ip1, err := parseStringToTime(t1)
		if err != nil {
			return nil, err
		}

		ip2, err := parseStringToTime(t2)
		if err != nil {
			return nil, err
		}

		loc, appErr := parseTimezone(tz)
		if appErr != nil {
			return nil, &ApplicationError{
				Message:    "cannot parse timezone",
				StatusCode: http.StatusBadRequest,
				Code:       "bad_request",
			}
		}

		return getTimestamps(ip1, ip2, loc, period)
	} else {
		return nil, &ApplicationError{
			Message:    "cannot parse invocation points",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}
}

func getTimestamps(ip1, ip2 *time.Time, loc *time.Location, period string) ([]string, *ApplicationError) {

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
		return getHourlyTimestamps(ip1, ip2, n), nil
	case "d":
		return getDailyTimestamps(ip1, ip2, loc, n), nil
	case "mo":
		return getMonthlyTimestamps(ip1, ip2, n), nil
	case "y":
		return getAnnuallyTimestamps(ip1, ip2, n), nil
	default:
		return nil, &ApplicationError{
			Message:    "could not parse period",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
	}
}

func getAnnuallyTimestamps(ip1, ip2 *time.Time, n int) []string {
	timestamps := []string{}
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(n, 0, 0)
	}
	return timestamps
}

func getMonthlyTimestamps(ip1, ip2 *time.Time, n int) []string {
	timestamps := []string{}
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, n, 0)
	}
	return timestamps
}

func getDailyTimestamps(ip1, ip2 *time.Time, loc *time.Location, n int) []string {
	timestamps := []string{}

	// Calculate the offset of the local zone. Offset is the actual difference between
	// UTC and Local time.
	_, offset := ip1.Local().Zone()

	timestamp := ip1.Truncate(time.Duration(offset) * time.Second)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Local().In(loc).Format(UTC_FORM))
		timestamp = timestamp.AddDate(0, 0, n)
	}
	return timestamps
}

func getHourlyTimestamps(ip1, ip2 *time.Time, n int) []string {
	timestamps := []string{}
	timestamp := ip1.Round(time.Hour)
	for timestamp.Before(*ip2) {
		timestamps = append(timestamps, timestamp.Format(UTC_FORM))
		timestamp = timestamp.Add(time.Duration(n) * time.Hour)
	}
	return timestamps
}
