package main

import (
	"os"

	"github.com/no-yan/wrkt/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
