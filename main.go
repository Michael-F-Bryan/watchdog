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

func check(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	up := Checksite(url)
	LogState(up, url)
}

func Checksite(url string) bool {

	respone, err := http.Get(url)
	if err != nil {
		return false
	}

	return respone.StatusCode == 200
}

func LogState(up bool, url string) {
	var state string
	if up {
		state = "up"
	} else {
		state = "down"
	}

	log.Printf("%v is %v", url, state)
}
