package apitest

import (
	"bytes"
	"fmt"
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
