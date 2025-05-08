package main

import (
	"os"

	cmd "github.com/gajirou/fexec"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
