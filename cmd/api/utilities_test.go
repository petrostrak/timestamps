package main

import (
	"testing"
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
