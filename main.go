package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
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
	Target    *WebTarget
	Timestamp time.Time
}

// WebTarget represents a service to be checked.
type WebTarget struct {
	Name    string
	URL     string
	Timeout time.Duration
}

type conf struct {
	WebTarget []struct {
		URL     string        `yaml:"cmt"`
		Timeout time.Duration `yaml:"con"`
	}
}

func (c *conf) getConfig() *conf {

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile Get err: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}
	return c
}

func main() {
	flag.Parse()

	var c conf
	//targets := c.getConfig()
	fmt.Println(c)

	// The URLs that will be checked when run
	targets := []WebTarget{
		{Name: "ConfLuence", URL: "http://134.7.57.175:8090/", Timeout: 10 * time.Second},
		{Name: "CMT", URL: "http://www.curtinmotorsport.com", Timeout: 10 * time.Second},
	}
	wg := sync.WaitGroup{}
	for _, target := range targets {
		wg.Add(1)
		go check(target, &wg)
	}
	wg.Wait()
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
