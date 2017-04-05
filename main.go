package main

import (
	"fmt"
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

// Status represents the state of a service at a particular point in time.
type Status struct {
	State     State
	Name      string
	Timestamp time.Time
}

func main() {
	// The urls that will be checked when run
	urls := []string{
		"http://134.7.57.175:8090/",
		"www.curtinmotorsport.com",
	}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go check(url, &wg)
	}
	wg.Wait()
}

// check will log a urls state and then tell the wait group that it is done.
func check(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	state := Checksite(url)
	LogState(state, url)
}

// Checksite will get a respone and error from a url and return the status of it.
// If any errors are encountered, expect 200, it's assumed that the site is down.
func Checksite(url string) Status {
	status := Status{Name: url, Timestamp: time.Now()}

	respone, err := http.Get(url)
	// TODO: inspect what error the page returns
	if err != nil {
		status.State = StateDown
		return status
	}
	if respone.StatusCode == 200 {
		status.State = StateUp
	} else {
		status.State = StateUnknown
	}

	return status
}

// LogState will log a struct of the state of the url, and other important info
// Just prints the struct but later on could save to a database.
func LogState(status Status, url string) {
	fmt.Printf("%#v\n", status)
}
