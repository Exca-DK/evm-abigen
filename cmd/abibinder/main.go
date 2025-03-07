package main

import (
	"fmt"
	"os"

	"github.com/Exca-DK/evm-abigen/cmd/abibinder/launcher"
)

func main() {
	exit(launcher.Launch(os.Args))
}

func exit(err error) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
