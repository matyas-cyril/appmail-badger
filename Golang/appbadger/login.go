package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	pkgConfig "appbadger/config"
	pkgLogin "appbadger/login"

	"github.com/matyas-cyril/argon2id"
	"golang.org/x/term"
)

func login(args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("command login is empty")
	}

	epoch := time.Now().Unix()

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			login, confFile, passwd, description string
			nbf, exp                             int64
			enable, passStdin                    bool
		)

		addFs := flag.NewFlagSet("add", flag.ExitOnError)
		addFs.StringVar(&login, "login", "", "username for login")
		addFs.BoolVar(&passStdin, "P", false, "input password")

		addFs.Int64Var(&nbf, "nbf", epoch, "not before")
		addFs.Int64Var(&exp, "exp", 0, "expiration")
		addFs.BoolVar(&enable, "enable", false, "enable")
		addFs.StringVar(&login, "description", "", "description")
		addFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")

		if err := addFs.Parse(args[1:]); err != nil {
			return err
		}

		if len(addFs.Args()) > 0 {
			return fmt.Errorf("unexpected arg : %q", addFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return err
		}

		if nbf < epoch {
			return fmt.Errorf("nbf must be greater than or equal to %d", epoch)
		}

		if exp != 0 && exp <= nbf {
			return fmt.Errorf("exp must be greater than or equal to nbf -> %d", epoch)
		}

		login = strings.TrimSpace(login)

		// Vérification des args obligatoires
		if login == "" || !passStdin {
			addFs.Usage()
			return fmt.Errorf("error: -login and -P are required at least")
		}

		if err := pkgLogin.IsValidLogin(login); err != nil {
			return err
		}

		if passStdin {
			var err error
			passwd, err = password()
			if err != nil {
				return err
			}
		}

		// Générer le mot de passe
		pwdByte, err := argon2id.Hash(passwd, argon2id.DefaultParams())
		if err != nil {
			return err
		}

		return pkgLogin.AddLogin(login, string(pwdByte), strings.TrimSpace(description), nbf, exp, !enable, *cfg)

	case "del":

		var (
			login, confFile, uuid string
		)

		delFs := flag.NewFlagSet("add", flag.ExitOnError)
		delFs.StringVar(&login, "login", "", "username for login")
		delFs.StringVar(&uuid, "uuid", "", "UUID for login")
		delFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")

		if err := delFs.Parse(args[1:]); err != nil {
			return err
		}

		if len(delFs.Args()) > 0 {
			return fmt.Errorf("unexpected arg : %q", delFs.Args()[0])
		}

		// Vérification des args
		if (login == "" && uuid == "") || (login != "" && uuid != "") {
			delFs.Usage()
			return fmt.Errorf("error: -login OR -uuid must defined")
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return err
		}

		if uuid != "" {
			return pkgLogin.DelLoginByUUID(uuid, *cfg)
		}

		login = strings.TrimSpace(login)
		if err := pkgLogin.IsValidLogin(login); err != nil {
			return err
		}

		return pkgLogin.DelLogin(login, *cfg)

	case "get":

		var (
			login, confFile, uuid string
			all                   bool
		)

		getFs := flag.NewFlagSet("get", flag.ExitOnError)
		getFs.StringVar(&login, "login", "", "username for login")
		getFs.StringVar(&uuid, "uuid", "", "uuid for login")
		getFs.BoolVar(&all, "all", false, "get all logins")
		getFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")

		if err := getFs.Parse(args[1:]); err != nil {
			return err
		}

		if len(getFs.Args()) > 0 {
			return fmt.Errorf("unexpected arg : %q", getFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return err
		}

		if all {
			return pkgLogin.GetLogins(*cfg)
		}

		if login != "" {
			return pkgLogin.GetLogin(login, *cfg)
		}

		if uuid != "" {
			return pkgLogin.GetLoginByUUID(uuid, *cfg)
		}

		getFs.Usage()
		return fmt.Errorf("error: -all  or -login or -uuid must defined")

	case "check": // Vérifier l'authentification

		var (
			login, passwd, confFile string
		)

		checkFs := flag.NewFlagSet("add", flag.ExitOnError)
		checkFs.StringVar(&login, "login", "", "username for login")
		checkFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")

		if err := checkFs.Parse(args[1:]); err != nil {
			return err
		}

		if len(checkFs.Args()) > 0 {
			return fmt.Errorf("unexpected arg : %q", checkFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return err
		}

		login = strings.TrimSpace(login)
		if err := pkgLogin.IsValidLogin(login); err != nil {
			return err
		}

		// Saisir le mot de passe
		fmt.Print("Password : ")
		passwd, ctrl, err := readPwdWithMask()
		if err != nil {
			return fmt.Errorf("input password error : %s", err.Error())
		}
		fmt.Println("")
		if ctrl {
			os.Exit(ERR_ARGS)
		}

		return pkgLogin.Check(login, passwd, *cfg)

	case "enable", "disable":

	default:
		fmt.Fprintf(os.Stdout, "command %s not exits\n", arg)

	}

	return nil
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
	if err := pkgLogin.IsValidPassword(pwd); err != nil {
		return "", err
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
		//if r > 32 && r < 127 {
		if pkgLogin.IsValidPasswordRune(r) {
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
