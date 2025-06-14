package main

import (
	"fmt"
)

func main() {
	cmd := newRootCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
