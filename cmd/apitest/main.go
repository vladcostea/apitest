package main

import (
	"github.com/vladcostea/apitest"
)

func main() {
	app := apitest.App{}
	app.Run(apitest.NewCLI())
}
