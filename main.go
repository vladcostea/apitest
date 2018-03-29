package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

type TestCase struct {
	Desc    string
	Path    string
	Json    string
	Status  int
	Headers []Header
}

type TestSuiteConfig struct {
	BaseURI string `yaml:"base_uri"`
	Auth    string
	Headers []Header
}

type TestSuite struct {
	TestSuiteConfig `yaml:"config"`
	Tests           []TestCase
}

type TestResult struct {
	Desc   string
	Error  string
	Passed bool
}

func loadTestSuite(filename string) TestSuite {
	yamlFile, err := ioutil.ReadFile(filename)
	s := TestSuite{}

	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &s)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return s
}

func (c TestSuiteConfig) HasAuth() bool {
	return len(c.Auth) > 0
}

func (c TestSuiteConfig) BasicAuth() (string, string) {
	auth := strings.Split(c.Auth, ":")
	return auth[0], auth[1]
}

func (c TestSuiteConfig) URL(path string) string {
	return fmt.Sprintf("%v%v", c.BaseURI, path)
}

func (t TestCase) Run(c TestSuiteConfig, results chan TestResult) {
	req, err := http.NewRequest("POST", c.URL(t.Path), bytes.NewBuffer([]byte(t.Json)))
	if err != nil {
		panic(err)
	}
	for _, h := range c.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	for _, h := range t.Headers {
		req.Header.Set(h.Name, h.Value)
	}
	if c.HasAuth() {
		user, pass := c.BasicAuth()
		req.SetBasicAuth(user, pass)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != t.Status {
		body, _ := ioutil.ReadAll(resp.Body)
		errorMessage := fmt.Sprintf("[FAIL] wanted %v got %v %v", t.Status, resp.StatusCode, string(body))
		results <- TestResult{
			Passed: false,
			Desc:   t.Desc,
			Error:  errorMessage,
		}
	} else {
		results <- TestResult{
			Passed: true,
		}
	}
}

func main() {
	suite := loadTestSuite("suite.yml")
	results := make(chan TestResult)
	for _, test := range suite.Tests {
		go test.Run(suite.TestSuiteConfig, results)
	}

	var failures []TestResult
	var total int
	for result := range results {
		total = total + 1
		if result.Passed {
			fmt.Print(".")
		} else {
			fmt.Print("F")
			failures = append(failures, result)
		}

		if total == len(suite.Tests) {
			close(results)
		}
	}

	if len(failures) > 0 {
		for index, failure := range failures {
			fmt.Printf("\n%v) %v", index+1, failure.Desc)
			fmt.Printf("\n\t%v\n", failure.Error)
		}
	} else {
		fmt.Printf("\nAll tests passed.\n")
	}
}
