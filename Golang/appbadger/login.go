package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	pkgLogin "appbadger/login"

	"golang.org/x/term"
)

func login(args []string) (*pkgLogin.LOGIN, error) {

	if len(args) == 0 {
		return nil, fmt.Errorf("command login is empty")
	}

	epoch := time.Now().Unix()

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			login, file, passwd string
			nbf, exp            int64
			enable, passStdin   bool
		)

		addFs := flag.NewFlagSet("add", flag.ExitOnError)
		addFs.StringVar(&login, "name", "", "username for login")
		addFs.BoolVar(&passStdin, "P", false, "input password")

		addFs.Int64Var(&nbf, "nbf", epoch, "not before")
		addFs.Int64Var(&exp, "exp", 0, "expiration")
		addFs.BoolVar(&enable, "enable", false, "enable")
		addFs.StringVar(&file, "file", DEFAULT_PATH, "file db")

		if err := addFs.Parse(args[1:]); err != nil {
			return nil, err
		}

		if len(addFs.Args()) > 0 {
			return nil, fmt.Errorf("unexpected arg : %q", addFs.Args()[0])
		}

		if nbf < epoch {
			return nil, fmt.Errorf("nbf must be greater than or equal to %d", epoch)
		}

		if exp != 0 && exp <= nbf {
			return nil, fmt.Errorf("exp must be greater than or equal to nbf -> %d", epoch)
		}

		login = strings.TrimSpace(login)

		// Vérification des args obligatoires
		if strings.TrimSpace(login) == "" || !passStdin {
			addFs.Usage()
			return nil, fmt.Errorf("error: -user and -P are required")
		}

		if passStdin {
			var err error
			passwd, err = password()
			if err != nil {
				return nil, err
			}
		}

		return &pkgLogin.LOGIN{
			Cmd:      "ADD",
			Login:    login,
			Password: passwd,
			Nbf:      nbf,
			Exp:      exp,
			Enable:   enable,
			File:     file,
		}, nil

	case "del":

	case "info":

	case "check": // Vérifier l'authentification

	case "enable", "disable":

	default:
		fmt.Fprintf(os.Stdout, "command %s not exits\n", arg)

	}

	return nil, nil
}

func password() (string, error) {

	// Traitement pour saisir le mot de passe
	fmt.Print("Password : ")
	pwd, ctrl, err := readPwdWithMask()
	if err != nil {
		return "", fmt.Errorf("input password error : %s\n", err.Error())
	}
	fmt.Println("")
	if ctrl {
		os.Exit(ERR_ARGS)
	}

	// Vérifier la taille min et max du mot de passe
	pwd = strings.TrimSpace(pwd)
	if len(pwd) < int(PASSWD_MIN_SIZE) || len(pwd) > int(PASSWD_MAX_SIZE) {
		return "", fmt.Errorf("password length invalid - valid lenght size min: %d - max: %d characters", PASSWD_MIN_SIZE, PASSWD_MAX_SIZE)
	}

	fmt.Print("Confirm  : ")
	confirmPwd, ctrl, err := readPwdWithMask()
	fmt.Println("")
	if ctrl {
		os.Exit(ERR_ARGS)
	}
	if err != nil {
		fmt.Fprintf(os.Stdout, "confirm password error : %s\n", err.Error())
		os.Exit(ERR_ARGS)
	}

	if strings.Compare(pwd, confirmPwd) != 0 {
		fmt.Println("Password mismatch")
	}

	return pwd, nil

}

func readPwdWithMask() (string, bool, error) {

	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return "", false, err
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	reader := bufio.NewReader(os.Stdin)
	var password string

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", false, err
		}

		// CTRL+C, Ctrl+D on quitte
		if r == 3 || r == 4 {
			fmt.Println("Quit")
			return "", true, nil
		}

		// Caractères que l'on refuse
		if r == '\t' || r == 27 {
			continue
		}

		if r == '\n' || r == '\r' {
			break
		}

		// Caractères que l'on va autoriser à sasir
		if r > 32 && r < 127 {

			password += string(r)
			fmt.Print("*")

		} else if r == '\b' || r == 127 {

			if len(password) > 0 {
				password = password[:len(password)-1]
				fmt.Print("\b \b")
			}

		} else {
			continue
		}

	}

	return password, false, nil
}
