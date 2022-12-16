package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTimestamps(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		queryParams        map[string]string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{"get all timestamps", "GET", map[string]string{"period": "1h", "tz": "Europe/Athens", "t1": "20210714T204603Z", "t2": "20210715T123456Z"}, getAllTimestamps, http.StatusOK},
		{"get all timestamps - bad t1", "GET", map[string]string{"period": "1h", "tz": "Europe/Athens", "t1": "", "t2": "20210715T123456Z"}, getAllTimestamps, http.StatusBadRequest},
		{"get all timestamps - bad timezone", "GET", map[string]string{"period": "1h", "tz": "Eur/Athens", "t1": "20210714T204603Z", "t2": "20210715T123456Z"}, getAllTimestamps, http.StatusBadRequest},
		{"get all timestamps - bad period", "GET", map[string]string{"period": "1w", "tz": "Europe/Athens", "t1": "20210714T204603Z", "t2": "20210715T123456Z"}, getAllTimestamps, http.StatusBadRequest},
	}

	for _, e := range testCases {

		req, _ := http.NewRequest(e.method, "/", nil)
		q := req.URL.Query()
		for k, v := range e.queryParams {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(e.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s: wrong status returned; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

	}
}
