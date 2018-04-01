package apitest

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type CLI struct {
	Total  int
	Suites []TestSuite
}

func NewCLI() *CLI {
	c := &CLI{}
	c.Load()

	return c
}

func (c *CLI) Load() {
	flag.Parse()
	files := flag.Args()
	var suites []TestSuite
	var totalTests int
	for _, filename := range files {
		suite, err := loadTestSuite(filename)
		if err != nil {
			fmt.Printf("test suite load error %v", err)
		} else {
			suites = append(suites, suite)
			totalTests += len(suite.Tests)
		}
	}

	c.Suites = suites
	c.Total = totalTests
}

type yamlSuite struct {
	Config struct {
		Base    string   `yaml:"base_uri"`
		Auth    string   `yaml:"auth"`
		Headers []header `yaml:"headers"`
	} `yaml:"config"`

	Tests []struct {
		Label   string   `yaml:"label"`
		Path    string   `yaml:"path"`
		Body    string   `yaml:"body"`
		Status  int      `yaml:"status"`
		Headers []header `yaml:"headers"`
	} `yaml:"tests"`
}

func loadTestSuite(filename string) (TestSuite, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return TestSuite{}, err
	}

	suite := yamlSuite{}
	err = yaml.Unmarshal(yamlFile, &suite)
	if err != nil {
		return TestSuite{}, err
	}

	var tests []TestCase
	for _, t := range suite.Tests {
		var headers []header
		for _, h := range suite.Config.Headers {
			headers = append(headers, h)
		}
		for _, h := range t.Headers {
			headers = append(headers, h)
		}

		url := fmt.Sprintf("%v%v", suite.Config.Base, t.Path)
		tests = append(tests, TestCase{
			Label:   t.Label,
			URL:     url,
			Body:    t.Body,
			Auth:    suite.Config.Auth,
			Headers: headers,
			Status:  t.Status,
		})
	}

	return TestSuite{Name: filename, Tests: tests}, nil
}
