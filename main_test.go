package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
func TestCheckWorkingSite(t *testing.T) {
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

		status := Checksite(server.URL)

		if status.State != input.ShouldBe {
			t.Errorf("Expected %v, got %v", input.ShouldBe, status.State)
		}
		server.Close()
	}
}
