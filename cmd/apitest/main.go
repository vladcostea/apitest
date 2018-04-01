package main

import (
	"fmt"
	"io/ioutil"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/vladcostea/apitest"
)

func loadSuite(filename string) (apitest.TestSuite, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return apitest.TestSuite{}, err
	}

	return apitest.NewTestSuiteFromYaml(filename, yamlFile)
}

func load(files []string) []apitest.TestSuite {
	var suites []apitest.TestSuite
	for _, filename := range files {
		suite, err := loadSuite(filename)
		if err != nil {
			fmt.Printf("[ERROR] test suite load error %v\n", err)
		} else {
			suites = append(suites, suite)
		}
	}

	return suites
}

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Printf("[ERROR] no input files provided\n")
		os.Exit(1)
	}

	tests := load(files)
	app := &apitest.App{}

	app.Run(tests)
}
