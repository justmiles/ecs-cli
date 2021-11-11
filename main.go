package main

import (
	"github.com/justmiles/ecs-cli/cmd"
)

// version of ecs-cli. Overwritten during build
var version = "development"

func main() {
	cmd.Execute(version)
}
