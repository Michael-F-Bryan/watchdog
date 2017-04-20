package main

import (
	"bytes"
	"flag"
	"io/ioutil"
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
var configFile = flag.String("config", "cfg.yaml", "The configuration file")

func main() {
	flag.Parse()

	targets, err := getConfig(*configFile)
	if err != nil {
		log.Fatalf("Error parsing the config file: %v", err)
	}

	wg := sync.WaitGroup{}
	for _, target := range targets {
		wg.Add(1)
		go check(target, &wg)
	}
	wg.Wait()
}

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

func getConfig(filename string) ([]WebTarget, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	targets, err := ParseConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	webTargets := make([]WebTarget, 0)
	for _, target := range targets {
		webTargets = append(webTargets, target.toTarget())
	}

	return webTargets, nil
}

// check will log a URLs state and then tell the wait group that it is done.
func check(target WebTarget, wg *sync.WaitGroup) {
	defer wg.Done()
	state := target.Check()
	LogState(state)
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

// LogState will log a struct of the state of the URL, and other important info
// Just prints the struct but later on could save to a database.
func LogState(status Status) {
	log.Printf("%v is %v", status.Target.URL, status.State)
}
