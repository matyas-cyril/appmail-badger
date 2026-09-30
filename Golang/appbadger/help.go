package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func printHelp() {
	fmt.Fprintf(os.Stdout, "Usage: %s [login|help]\n\n", filepath.Base(os.Args[0]))
}
