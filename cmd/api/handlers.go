package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// PeriodicTask struct is the collection of attributes for periodic task
type PeriodicTask struct {
	Period   string `json:"period"`
	Timezone string `json:"tz"`
	InvocationPoints
	Timestamps []string `json:"timestamps"`
}

type InvocationPoints struct {
	T1 string `json:"t1"`
	T2 string `json:"t2"`
}

func GetAllTimestamps(w http.ResponseWriter, r *http.Request) {
	// read query parameteres
	queryParams := r.URL.Query()

	period := queryParams.Get("period")
	timezone := queryParams.Get("tz")
	t1, t2 := queryParams.Get("t1"), queryParams.Get("t2")

	tz, err := parseTimezone(timezone)
	if err != nil {
		RespondError(w, err)
		return
	}

	invocationPoints, err := parseInvocationPoints(t1, t2, period)
	if err != nil {
		RespondError(w, err)
		return
	}

	pd := &PeriodicTask{
		Period:   period,
		Timezone: tz,
		InvocationPoints: InvocationPoints{
			T1: t1,
			T2: t2,
		},
		Timestamps: invocationPoints,
	}

	Respond(w, http.StatusOK, pd)
}

type respond struct {
	Status int
	Body   any
}

// Respond is the function that returns a successful JSON
func Respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &respond{
		Status: status,
		Body:   body,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatalf("JSON decoding failed: %s", err)
	}

}

// ApplicationError is the struct responsible for errors
type ApplicationError struct {
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
	Message    string `json:"desc"`
}

// RespondError is the function that returns an ApplicationError
func RespondError(w http.ResponseWriter, err *ApplicationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	resp := &respond{
		Status: err.StatusCode,
		Body:   err,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatalf("JSON decoding failed: %s", err)
	}
}
