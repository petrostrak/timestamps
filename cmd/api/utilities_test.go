package main

import (
	"testing"
	"time"
)

func TestIsValidPeriod(t *testing.T) {
	testCases := []struct {
		name     string
		period   string
		expected bool
	}{
		{"hourly", "1h", true},
		{"daily", "1d", true},
		{"monthly", "1mo", true},
		{"annually", "1y", true},
		{"weekly", "1w", false},
	}

	for _, e := range testCases {
		actual := isValidPeriod(e.period)

		if actual != e.expected {
			t.Errorf("Expected %v but got %v", e.expected, actual)
		}
	}
}

func TestParseTimeZone(t *testing.T) {
	testCases := []struct {
		name     string
		timezone string
		expected string
	}{
		{"correct", "Europe/Athens", "Europe/Athens"},
		{"false", "Eu/Ath", ""},
	}

	for _, e := range testCases {
		actual, _ := parseTimezone(e.timezone)

		if actual != e.expected {
			t.Errorf("Expected %s but got %s", e.expected, actual)
		}
	}
}

func TestCheckInvocationPoints(t *testing.T) {
	testCases := []struct {
		name          string
		invocation_p1 string
		invocation_p2 string
		expected      bool
	}{
		{"correct", "20210714T204603Z", "20210715T123456Z", true},
		{"wrong", "20210714204603", "20210715123456", false},
	}

	for _, e := range testCases {
		result := CheckInvocationPoints(e.invocation_p1, e.invocation_p2)

		if result != e.expected {
			t.Errorf("Expected %v but got %v", e.expected, result)
		}
	}
}

func TestCheckInvocationSequence(t *testing.T) {
	testCases := []struct {
		name          string
		invocation_p1 string
		invocation_p2 string
		layout        string
		expected      bool
	}{
		{"correct", "20210714T204603Z", "20210715T123456Z", "20060102T150405Z", true},
		{"wrong", "20210714T204603Z", "20210715T123456Z", time.RFC3339, false},
	}

	for _, e := range testCases {
		result := CheckInvocationSequence(e.invocation_p1, e.invocation_p2, e.layout)

		if result != e.expected {
			t.Errorf("Expected %v but got %v", e.expected, result)
		}
	}
}
