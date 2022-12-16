package main

import (
	"net/http"
	"reflect"
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

func TestParseStringToTime(t *testing.T) {
	testCases := []struct {
		name            string
		layout          string
		invocationPoint string
		expected        string
		err             *ApplicationError
	}{
		{"first correct", "20060102T150405Z", "20210714T204603Z", "2021-07-14 20:46:03 +0000 UTC", nil},
		{"second correct", "20060102T150405Z", "20210715T123456Z", "2021-07-15 12:34:56 +0000 UTC", nil},
		{"different layout", time.Layout, "20210715T123456Z", "", &ApplicationError{
			Message:    "cannot parse invocation points",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}},
	}

	for _, e := range testCases {
		result, err := ParseStringToTime(e.layout, e.invocationPoint)
		if err != nil {
			if err.Message != "cannot parse invocation points" {
				t.Error("Expected error, but got none")
			}
			return
		}

		if result.String() != e.expected {
			t.Errorf("%s: Expected %s but got %s", e.name, e.expected, result.String())
		}
	}
}

func TestGetHourlyTimestamps(t *testing.T) {
	testCases := []struct {
		name          string
		layout        string
		invocation_p1 string
		invocation_p2 string
		frequency     int
		timestamps    []string
		expected      []string
	}{
		{"correct", "20060102T150405Z", "20210714T204603Z", "20210715T123456Z", 1, []string{}, []string{"20210714T210000Z", "20210714T220000Z", "20210714T230000Z", "20210715T000000Z", "20210715T010000Z", "20210715T020000Z", "20210715T030000Z", "20210715T040000Z", "20210715T050000Z", "20210715T060000Z", "20210715T070000Z", "20210715T080000Z", "20210715T090000Z", "20210715T100000Z", "20210715T110000Z", "20210715T120000Z"}},
		{"correct - every 3 hours", "20060102T150405Z", "20210714T204603Z", "20210715T123456Z", 3, []string{}, []string{"20210714T210000Z", "20210715T000000Z", "20210715T030000Z", "20210715T060000Z", "20210715T090000Z", "20210715T120000Z"}},
		{"wrong invocation points", "20060102T150405Z", "20210715T123456Z", "20210714T204603Z", 1, []string{}, []string{}},
	}

	for _, e := range testCases {
		ip1, _ := ParseStringToTime(e.layout, e.invocation_p1)
		ip2, _ := ParseStringToTime(e.layout, e.invocation_p2)

		result := getHourlyTimestamps(ip1, ip2, e.timestamps, e.frequency)

		if !reflect.DeepEqual(result, e.expected) {
			t.Errorf("%s: Expected %v but got %v", e.name, e.expected, result)
		}
	}

}

func TestGetDailyTimestamps(t *testing.T) {
	testCases := []struct {
		name          string
		layout        string
		invocation_p1 string
		invocation_p2 string
		frequency     int
		timestamps    []string
		expected      []string
	}{
		{"correct", "20060102T150405Z", "20211010T204603Z", "20211115T123456Z", 1, []string{}, []string{"20211010T210000Z", "20211011T210000Z", "20211012T210000Z", "20211013T210000Z", "20211014T210000Z", "20211015T210000Z", "20211016T210000Z", "20211017T210000Z", "20211018T210000Z", "20211019T210000Z", "20211020T210000Z", "20211021T210000Z", "20211022T210000Z", "20211023T210000Z", "20211024T210000Z", "20211025T210000Z", "20211026T210000Z", "20211027T210000Z", "20211028T210000Z", "20211029T210000Z", "20211030T210000Z", "20211031T210000Z", "20211101T210000Z", "20211102T210000Z", "20211103T210000Z", "20211104T210000Z", "20211105T210000Z", "20211106T210000Z", "20211107T210000Z", "20211108T210000Z", "20211109T210000Z", "20211110T210000Z", "20211111T210000Z", "20211112T210000Z", "20211113T210000Z", "20211114T210000Z"}},
		{"correct - every 5 days", "20060102T150405Z", "20211010T204603Z", "20211115T123456Z", 5, []string{}, []string{"20211010T210000Z", "20211015T210000Z", "20211020T210000Z", "20211025T210000Z", "20211030T210000Z", "20211104T210000Z", "20211109T210000Z", "20211114T210000Z"}},
		{"wrong invocation points", "20060102T150405Z", "20211115T123456Z", "20211010T204603Z", 1, []string{}, []string{}},
	}

	for _, e := range testCases {
		ip1, _ := ParseStringToTime(e.layout, e.invocation_p1)
		ip2, _ := ParseStringToTime(e.layout, e.invocation_p2)

		result := getDailyTimestamps(ip1, ip2, e.timestamps, e.frequency)

		if !reflect.DeepEqual(result, e.expected) {
			t.Errorf("%s: Expected %v but got %v", e.name, e.expected, result)
		}
	}
}

func TestGetMonthlyTimestamps(t *testing.T) {
	testCases := []struct {
		name          string
		layout        string
		invocation_p1 string
		invocation_p2 string
		frequency     int
		timestamps    []string
		expected      []string
	}{
		{"correct", "20060102T150405Z", "20210214T204603Z", "20211115T123456Z", 1, []string{}, []string{"20210214T210000Z", "20210314T210000Z", "20210414T210000Z", "20210514T210000Z", "20210614T210000Z", "20210714T210000Z", "20210814T210000Z", "20210914T210000Z", "20211014T210000Z", "20211114T210000Z"}},
		{"correct - every 2 months", "20060102T150405Z", "20210214T204603Z", "20211115T123456Z", 2, []string{}, []string{"20210214T210000Z", "20210414T210000Z", "20210614T210000Z", "20210814T210000Z", "20211014T210000Z"}},
		{"wrong invocation points", "20060102T150405Z", "20211115T123456Z", "20210214T204603Z", 1, []string{}, []string{}},
	}

	for _, e := range testCases {
		ip1, _ := ParseStringToTime(e.layout, e.invocation_p1)
		ip2, _ := ParseStringToTime(e.layout, e.invocation_p2)

		result := getMonthlyTimestamps(ip1, ip2, e.timestamps, e.frequency)

		if !reflect.DeepEqual(result, e.expected) {
			t.Errorf("%s: Expected %v but got %v", e.name, e.expected, result)
		}
	}
}

func TestGetAnnuallyTimestamps(t *testing.T) {
	testCases := []struct {
		name          string
		layout        string
		invocation_p1 string
		invocation_p2 string
		frequency     int
		timestamps    []string
		expected      []string
	}{
		{"correct", "20060102T150405Z", "20180214T204603Z", "20211115T123456Z", 1, []string{}, []string{"20180214T210000Z", "20190214T210000Z", "20200214T210000Z", "20210214T210000Z"}},
		{"correct - every 2 years", "20060102T150405Z", "20180214T204603Z", "20211115T123456Z", 2, []string{}, []string{"20180214T210000Z", "20200214T210000Z"}},
		{"wrong invocation points", "20060102T150405Z", "20211115T123456Z", "20180214T204603Z", 1, []string{}, []string{}},
	}

	for _, e := range testCases {
		ip1, _ := ParseStringToTime(e.layout, e.invocation_p1)
		ip2, _ := ParseStringToTime(e.layout, e.invocation_p2)

		result := getAnnuallyTimestamps(ip1, ip2, e.timestamps, e.frequency)

		if !reflect.DeepEqual(result, e.expected) {
			t.Errorf("%s: Expected %v but got %v", e.name, e.expected, result)
		}
	}
}

func TestGetTimestamps(t *testing.T) {
	testCases := []struct {
		name          string
		layout        string
		invocation_p1 string
		invocation_p2 string
		period        string
		expected      []string
	}{
		{"hourly", "20060102T150405Z", "20210714T204603Z", "20210715T123456Z", "1h", []string{"20210714T210000Z", "20210714T220000Z", "20210714T230000Z", "20210715T000000Z", "20210715T010000Z", "20210715T020000Z", "20210715T030000Z", "20210715T040000Z", "20210715T050000Z", "20210715T060000Z", "20210715T070000Z", "20210715T080000Z", "20210715T090000Z", "20210715T100000Z", "20210715T110000Z", "20210715T120000Z"}},
		{"daily", "20060102T150405Z", "20211010T204603Z", "20211115T123456Z", "1d", []string{"20211010T210000Z", "20211011T210000Z", "20211012T210000Z", "20211013T210000Z", "20211014T210000Z", "20211015T210000Z", "20211016T210000Z", "20211017T210000Z", "20211018T210000Z", "20211019T210000Z", "20211020T210000Z", "20211021T210000Z", "20211022T210000Z", "20211023T210000Z", "20211024T210000Z", "20211025T210000Z", "20211026T210000Z", "20211027T210000Z", "20211028T210000Z", "20211029T210000Z", "20211030T210000Z", "20211031T210000Z", "20211101T210000Z", "20211102T210000Z", "20211103T210000Z", "20211104T210000Z", "20211105T210000Z", "20211106T210000Z", "20211107T210000Z", "20211108T210000Z", "20211109T210000Z", "20211110T210000Z", "20211111T210000Z", "20211112T210000Z", "20211113T210000Z", "20211114T210000Z"}},
		{"monthly", "20060102T150405Z", "20210214T204603Z", "20211115T123456Z", "1mo", []string{"20210214T210000Z", "20210314T210000Z", "20210414T210000Z", "20210514T210000Z", "20210614T210000Z", "20210714T210000Z", "20210814T210000Z", "20210914T210000Z", "20211014T210000Z", "20211114T210000Z"}},
		{"annually", "20060102T150405Z", "20180214T204603Z", "20211115T123456Z", "1y", []string{"20180214T210000Z", "20190214T210000Z", "20200214T210000Z", "20210214T210000Z"}},
		{"default", "20060102T150405Z", "20180214T204603Z", "20211115T123456Z", "false_period", []string{}},
	}

	for _, e := range testCases {
		ip1, _ := ParseStringToTime(e.layout, e.invocation_p1)
		ip2, _ := ParseStringToTime(e.layout, e.invocation_p2)

		result, err := GetTimestamps(ip1, ip2, e.period)
		if err != nil {
			if err.Message != "could not parse period" {
				t.Errorf("%s: Expected error but got none.", e.name)
			}
			return
		}

		if !reflect.DeepEqual(result, e.expected) {
			t.Errorf("%s: Expected %v but got %v", e.name, e.expected, result)
		}
	}
}
