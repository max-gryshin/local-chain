package main

import (
	"fmt"
	"os"

	"local-chain/tools/debug"
)

func main() {
	debugApp := debug.NewDebug()

	if err := debugApp.CMD.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
