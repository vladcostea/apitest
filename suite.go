package apitest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type header struct {
	Name  string
	Value string
}

type TestCase struct {
	Label   string
	URL     string
	Body    string
	Auth    string
	Status  int
	Headers []header
}

type TestSuite struct {
	Name  string
	Tests []TestCase
}

type TestResult struct {
	Desc   string
	Error  string
	Passed bool
}

func (c TestCase) HasAuth() bool {
	return len(c.Auth) > 0
}

func (c TestCase) BasicAuth() (string, string) {
	auth := strings.Split(c.Auth, ":")
	return auth[0], auth[1]
}

func (t TestCase) Run(results chan TestResult) {
	req, err := http.NewRequest("POST", t.URL, bytes.NewBuffer([]byte(t.Body)))
	if err != nil {
		panic(err)
	}

	for _, h := range t.Headers {
		req.Header.Set(h.Name, h.Value)
	}

	if t.HasAuth() {
		req.SetBasicAuth(t.BasicAuth())
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
			Desc:   t.Label,
			Error:  errorMessage,
		}
	} else {
		results <- TestResult{
			Passed: true,
		}
	}
}
