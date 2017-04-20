package main

import (
	"net/http"
	"time"
)

// State is an enum which represents the current state of a resource.
type State int

func (s State) String() string {
	switch s {
	case StateUp:
		return "up"
	case StateDown:
		return "down"
	case StateUnknown:
		return "unknown"
	default:
		panic("Unknown state")
	}
}

const (
	// StateUp represents a service which is currently up and responding.
	StateUp = iota

	// StateDown is a service which isn't responding.
	StateDown

	// StateUnknown means a service is in an unknown state.
	StateUnknown
)

// Status represents the state of a service at a particular point in time.
type Status struct {
	State     State
	Target    *WebTarget
	Timestamp time.Time
}

// WebTarget represents a service to be checked.
type WebTarget struct {
	Name    string
	URL     string
	Timeout time.Duration
}

// Checksite will get a response and error from a URL and return the status of it.
// If any errors are encountered, expect 200, it's assumed that the site is down.
func (w *WebTarget) Check() Status {
	status := Status{Target: w, Timestamp: time.Now()}

	client := http.Client{
		Timeout: *timeout,
	}

	response, err := client.Get(w.URL)
	// TODO: inspect what error the page returns
	if err != nil {
		status.State = StateDown
	} else if response.StatusCode == 200 {
		status.State = StateUp
	} else {
		status.State = StateDown
	}

	return status
}
