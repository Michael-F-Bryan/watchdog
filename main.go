package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
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
	up := Checksite(url)
	LogState(up, url)
}

// Checksite will get a respone and error from a url and return the status of it.
// If any errors and encounterded it is assumed that the site is down.
func Checksite(url string) bool {
	respone, err := http.Get(url)
	// TODO: inspect what error the page returns
	if err != nil {
		return false
	}
	return respone.StatusCode == 200
}

// LogState will log the state of the given url, for now jsut prints to the screen
// Later on could save to a database.
func LogState(up bool, url string) {
	var state string
	if up {
		state = "up"
	} else {
		state = "down"
	}
	log.Printf("%v is %v", url, state)
}
