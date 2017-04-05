package main

import (
	"net/http"
	"sync"
    "fmt"
    "time"
)

type status struct {
    state string
    time time.Time
}

func main() {
    // The urls that will be checked when run
	urls := []string {
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
func Checksite(url string) (string) {
	respone, err := http.Get(url)
	// TODO: inspect what error the page returns
	if err != nil {
		return "Down"
	}
    if respone.StatusCode == 200 {
        return "Up"
    } else {
        return "Unknown"
	}
}

// LogState will log a struct of the state of the url, and other important info
// Just prints the struct but later on could save to a database.
func LogState(state string, url string) {
    now := time.Now()
    fmt.Println(url, status{state: state, time: now})
}
