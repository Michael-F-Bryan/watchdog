package main

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

type DodgyReader struct{}

func (_ DodgyReader) Read(_ []byte) (int, error) {
	return 0, errors.New("Oops...")
}

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
		if *targets[i] != shouldBe[i] {
			t.Errorf("Got %v, but expected %v", targets, shouldBe)
		}
	}
}

func TestConfigIsInvalidIfRequiredFieldsAreMissing(t *testing.T) {
	inputs := []string{
		`
- foo: Bar
  timeout: 7
`,
		`
- name: Confluence
  timeout: 5
`,
		`
- url: http://www.google.com/
  timeout: 5
`,
	}

	for _, input := range inputs {
		targets, err := ParseConfig(strings.NewReader(input))
		if err == nil || targets != nil {
			t.Errorf("Expected an error, got %v and %v", targets, err)
		}
	}
}

func TestTimeoutIsSetToADefault(t *testing.T) {
	src := `
- name: Confluence
  url: http://google.com/
`

	targets, err := ParseConfig(strings.NewReader(src))
	if err != nil {
		t.Errorf("Unexpected error, %v", err)
	}

	if len(targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(targets))
	}

	if targets[0].Timeout != 5 {
		t.Errorf("Expected a timeout of 5, got %v", targets[0].Timeout)
	}
}

func TestGetConfigConvenienceFunction(t *testing.T) {
	filename := "cfg.yaml"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Default config file not found!")
	}

	targets, err := getConfig(filename)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	shouldBe := []WebTarget{
		{Name: "Confluence", URL: "http://134.7.57.175:8090/", Timeout: 5 * time.Second},
		{Name: "Website", URL: "http://curtinmotorsport.com/", Timeout: 10 * time.Second},
	}

	if len(targets) != len(shouldBe) {
		t.Fatalf("Expected %d targets, got %d", len(shouldBe), len(targets))
	}

	for i := range targets {
		if targets[i] != shouldBe[i] {
			t.Errorf("Expected %#v, got %#v", targets[i], shouldBe[i])
		}
	}
}

func TestConfigFileParseError(t *testing.T) {
	targets, err := ParseConfig(strings.NewReader("\t"))
	if err == nil || targets != nil {
		t.Errorf("Expected an error, got %v and %#v", err, targets)
	}
}

func TestConfigFileDetectsReadErrors(t *testing.T) {
	targets, err := ParseConfig(DodgyReader{})
	if err == nil || targets != nil {
		t.Errorf("Expected an error, got %v and %#v", err, targets)
	}
}
