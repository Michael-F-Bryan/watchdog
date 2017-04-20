package main

import (
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/go-yaml/yaml"
)

type tempTarget struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

func (t tempTarget) toTarget() WebTarget {
	return WebTarget{
		Name:    t.Name,
		URL:     t.URL,
		Timeout: time.Duration(t.Timeout) * time.Second,
	}
}

func ParseConfig(reader io.Reader) ([]*tempTarget, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var targets []*tempTarget

	err = yaml.Unmarshal(data, &targets)
	if err != nil {
		return nil, err
	}

	// Go through the targets and make sure they all have a name and URL
	for _, target := range targets {
		if target.Name == "" || target.URL == "" {
			return nil, errors.New("Missing required fields")
		}

		// default timeout
		if target.Timeout == 0 {
			target.Timeout = 5
		}
	}

	return targets, nil
}
