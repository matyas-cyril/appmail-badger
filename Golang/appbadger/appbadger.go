package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	BUILD_DATE    string
	BUILD_TIME    string
	BUILD_VERSION string
)

const (
	OK int = iota
	ERR_ARGS
	ERR_CONF_FILE
	ERR_TMPL
	ERR_LOGIN
	ERR_PASSWORD
	ERR_EPOCH

	DEFAULT_CONFIG_FILE string = "appbadger.json"
)

func main() {

	initBuildVar()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Use: %s [user|date|help] {...options}\n", filepath.Base(os.Args[0]))
		os.Exit(ERR_ARGS)
	}

	switch arg := strings.TrimSpace(os.Args[1]); arg {

	case "user":

		codeRetour, err := user(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", err)
		}
		os.Exit(codeRetour)

	case "date":
		codeRetour, err := date(os.Args[2:])
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

func initBuildVar() {
	if strings.TrimSpace(BUILD_DATE) == "" {
		BUILD_DATE = time.Now().Format("2006-01-02")
	}

	if strings.TrimSpace(BUILD_TIME) == "" {
		BUILD_TIME = time.Now().Format("15:04:05")
	}

	if strings.TrimSpace(BUILD_VERSION) == "" {
		BUILD_VERSION = "0.0.0"
	}

}
