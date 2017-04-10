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
	server := getHTTPServer(http.StatusOK)
	defer server.Close()

	status := Checksite(server.URL)

	if status.State != StateUp {
		t.Errorf("Expected StateUp, got %v", status.State)
	}
}

func TestCheckDownSite(t *testing.T) {
	server := getHTTPServer(http.StatusNotFound)
	defer server.Close()

	status := Checksite(server.URL)

	if status.State != StateDown {
		t.Errorf("Expected StateDown, got %v", status.State)
    }
}
