package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type TestCase struct {
	Desc   string
	Path   string
	Json   string
	Status int
}

type Config struct {
	BaseURI string `yaml:"base_uri"`
	Tests   []TestCase
}

type TestResult struct {
	Desc   string
	Error  string
	Passed bool
}

func LoadConfig(filename string) Config {
	yamlFile, err := ioutil.ReadFile(filename)
	c := Config{}

	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return c
}

func buildUrl(base, path string) string {
	return fmt.Sprintf("%v%v", base, path)
}

func (t TestCase) Run(baseURI string, results chan TestResult) {
	req, err := http.NewRequest("POST", buildUrl(baseURI, t.Path), bytes.NewBuffer([]byte(t.Json)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/pdf")
	req.SetBasicAuth(os.Getenv("USER"), os.Getenv("PASS"))

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
	config := LoadConfig("config.yml")
	results := make(chan TestResult)
	for _, test := range config.Tests {
		go test.Run(config.BaseURI, results)
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

		if total == len(config.Tests) {
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
