package main

import (
	"flag"
	"log"
	"net/http"
	"sync"
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

// Time before time out
var timeout = flag.Duration("timeout", 10*time.Second, "Timeout in seconds")

// Status represents the state of a service at a particular point in time.
type Status struct {
	State     State
	Name      string
	Timestamp time.Time
}

type WebTarget struct {
    url string
    timeout time.Duration
}

func main() {
	flag.Parse()

	// The urls that will be checked when run
	targets := []WebTarget{
		{"http://134.7.57.175:8090/", 10 * time.Second},
		{"http://www.curtinmotorsport.com", 10 * time.Second},
	}
	wg := sync.WaitGroup{}
	for _, target := range targets {
		wg.Add(1)
		go check(target, &wg)
	}
	wg.Wait()
}

// check will log a urls state and then tell the wait group that it is done.
func check(target WebTarget, wg *sync.WaitGroup) {
	defer wg.Done()
	state := target.Check()
	LogState(state, target.url)
}

// Checksite will get a response and error from a url and return the status of it.
// If any errors are encountered, expect 200, it's assumed that the site is down.
func (w *WebTarget) Check() Status {
	status := Status{Name: w.url, Timestamp: time.Now()}

	client := http.Client{
		Timeout: *timeout,
	}

	response, err := client.Get(w.url)
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

// LogState will log a struct of the state of the url, and other important info
// Just prints the struct but later on could save to a database.
func LogState(status Status, url string) {
	log.Printf("%v is %v", status.Name, status.State)
}
