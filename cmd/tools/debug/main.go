package main

import (
	"fmt"
	"os"

	debug2 "local-chain/internal/pkg/debug"
)

func main() {
	debugApp := debug2.NewDebug()

	if err := debugApp.CMD.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
