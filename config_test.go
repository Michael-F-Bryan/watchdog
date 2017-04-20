package main

import (
	"strings"
	"testing"
)

func TestReading(t *testing.T) {
	src := `
        - name: Confluence
          url: http://134.7.57.175:8090/
          timeout: 5
        `

	shouldBe := []tempTarget{
		{Name: "Confluence", URL: "http://134.7.57.175:8090/", Timeout: 5},
	}

	targets, err := ParseConfig(strings.NewReader(src))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(targets) != len(shouldBe) {
		t.Fatalf("Expected %d elements, got %d", len(shouldBe), len(targets))
	}

	for i := range targets {
		if targets[i] != shouldBe[i] {
			t.Errorf("Got %v, but expected %v", targets, shouldBe)
		}
	}
}
