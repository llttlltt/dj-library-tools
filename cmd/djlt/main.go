package main

import (
	"fmt"
	"os"

	"github.com/llttlltt/dj-library-tools/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
