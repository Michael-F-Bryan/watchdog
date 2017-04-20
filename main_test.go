package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func getHTTPServer(statusCode int) *httptest.Server {
	callback := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}
	server := httptest.NewServer(http.HandlerFunc(callback))

	return server
}

// TestCheckWorkingSite starts up a web server on the local machine then makes
// sure the `Checksite` function shows it's up
func TestSiteCheck(t *testing.T) {
	inputs := []struct {
		Status   int
		ShouldBe State
	}{
		{http.StatusOK, StateUp},
		{http.StatusNotFound, StateDown},
	}

	for _, input := range inputs {
		t.Log(input)
		server := getHTTPServer(input.Status)
		target := WebTarget{URL: server.URL, Timeout: 1 * time.Second}

		status := target.Check()

		if status.State != input.ShouldBe {
			t.Errorf("Expected %v, but got %v", input.ShouldBe, status.State)
		}
		server.Close()
	}
}

// TestDown starts up a web server on the local machine then makes sure the
// `Checksite` function shows that it's down when given an empty GET request
func TestDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {}))

	target := WebTarget{URL: server.URL, Timeout: 1 * time.Second}

	server.Close()
	Status := target.Check()

	if Status.State != StateDown {
		t.Errorf("Expected StateDown, but got %v", Status.State)
	}
}

// TestState checks the 3 outcomes of (s State), expect for default
func TestState(t *testing.T) {
	inputs := []struct {
		State    State
		ShouldBe string
	}{
		{StateUp, "up"},
		{StateDown, "down"},
		{StateUnknown, "unknown"},
	}

	for _, input := range inputs {
		t.Log(input)
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		stateStr := State.String(input.State)

		if stateStr != input.ShouldBe {
			t.Errorf("Expect %v, but got %v", input.ShouldBe, stateStr)
		}
	}
}
