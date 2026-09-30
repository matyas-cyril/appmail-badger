package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	OK int = iota
	ERR_LOGIN
	ERR_ARGS

	DEFAULT_PATH string = "/etc/file.js"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "Use: %s [login|help] {...options}\n", filepath.Base(os.Args[0]))
		os.Exit(OK)
	}

	//epoch := time.Now().Unix()

	switch arg := strings.TrimSpace(os.Args[1]); arg {

	case "login":
		err := login(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(ERR_ARGS)
		}

	case "help":
		printHelp()

	default:
		fmt.Fprintf(os.Stdout, "command %s not exits\n", arg)

	}

}
