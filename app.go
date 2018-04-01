package apitest

import (
	"fmt"
)

type App struct{}

func totalTests(suites []TestSuite) (total int) {
	for _, suite := range suites {
		total += len(suite.Tests)
	}

	return
}

func (app *App) Run(suites []TestSuite) {
	results := make(chan TestResult)
	var failures []TestResult
	var testsRan int
	total := totalTests(suites)

	for _, s := range suites {
		for _, t := range s.Tests {
			go t.Run(results)
		}
	}

	for result := range results {
		testsRan = testsRan + 1
		if result.Passed {
			fmt.Print(".")
		} else {
			fmt.Print("F")
			failures = append(failures, result)
		}

		if testsRan == total {
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
