package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	OK int = iota
	ERR_LOGIN
	ERR_ARGS

	DEFAULT_CONFIG_FILE string = "appbadger.json"
)

type ErrorResponse struct {
	Err string `json:"err"`
}

func stdoutErrToJSON(err error) {

	if err == nil {
		return
	}

	// Écriture du json sur stdout
	json.NewEncoder(os.Stdout).Encode(ErrorResponse{
		Err: err.Error(),
	})
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Use: %s [login|help] {...options}\n", filepath.Base(os.Args[0]))
		os.Exit(ERR_ARGS)
	}

	switch arg := strings.TrimSpace(os.Args[1]); arg {

	case "login":

		if err := login(os.Args[2:]); err != nil {
			stdoutErrToJSON(err)
			os.Exit(ERR_LOGIN)
		}

	case "help":
		printHelp("full_help")

	default:
		stdoutErrToJSON(fmt.Errorf("command %s not exist", arg))
		os.Exit(ERR_ARGS)

	}
	os.Exit(OK)
}
