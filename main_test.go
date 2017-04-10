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
		target := WebTarget{server.URL, 1 * time.Second}

		status := target.Check()

		if status.State != input.ShouldBe {
			t.Errorf("Expected %v, got %v", input.ShouldBe, status.State)
		}
		server.Close()
	}
}

func TestDown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {}))

	target := WebTarget{server.URL, 1 * time.Millisecond}

	server.Close()
	Status := target.Check()

	if Status.State != StateDown {
		t.Errorf("Expected StateDown, got %v", Status.State)
	}
}
