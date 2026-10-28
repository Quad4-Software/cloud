package main

import (
	"fmt"
	"os"

	"github.com/Quad4-Software/cloud/quad4/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
