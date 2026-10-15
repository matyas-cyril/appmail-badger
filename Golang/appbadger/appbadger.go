package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	OK int = iota
	ERR_ARGS
	ERR_CONF_FILE
	ERR_TMPL
	ERR_LOGIN
	ERR_PASSWORD

	DEFAULT_CONFIG_FILE string = "appbadger.json"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Use: %s [user|help] {...options}\n", filepath.Base(os.Args[0]))
		os.Exit(ERR_ARGS)
	}

	switch arg := strings.TrimSpace(os.Args[1]); arg {

	case "user":

		codeRetour, err := user(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
		}
		os.Exit(codeRetour)

	case "help":
		if err := printHelp("full_help"); err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
			os.Exit(ERR_TMPL)
		}

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exist\n", arg)
		os.Exit(ERR_ARGS)

	}
	os.Exit(OK)
}
