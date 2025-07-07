package main

import (
	"os"

	"github.com/no-yan/wt/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
