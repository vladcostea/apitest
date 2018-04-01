package apitest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type App struct{}

func (app App) Run() {
	suite := loadTestSuite("suite.yml")
	results := make(chan TestResult)
	for _, t := range suite.Tests {
		go t.Run(suite.TestSuiteConfig, results)
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
