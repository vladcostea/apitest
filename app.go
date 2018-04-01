package apitest

import (
	"fmt"
)

type App struct{}

func (app App) Run(c *CLI) {
	results := make(chan TestResult)
	var failures []TestResult
	var total int

	for _, suite := range c.Suites {
		for _, t := range suite.Tests {
			go t.Run(results)
		}
	}

	for result := range results {
		total = total + 1
		if result.Passed {
			fmt.Print(".")
		} else {
			fmt.Print("F")
			failures = append(failures, result)
		}

		if total == c.Total {
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
