package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"sync"
	"time"
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

// LogState will log a struct of the state of the URL, and other important info
// Just prints the struct but later on could save to a database.
func LogState(status Status) {
	log.Printf("%v is %v", status.Target.URL, status.State)
}
