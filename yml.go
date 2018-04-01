package apitest

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

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

func NewTestSuiteFromYaml(name string, yamlFile []byte) (TestSuite, error) {
	suite := yamlSuite{}
	err := yaml.Unmarshal(yamlFile, &suite)
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

	return TestSuite{Name: name, Tests: tests}, nil
}
