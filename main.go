package main

import (
	"os"
)

func main() {
	cmd := newRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
