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

func login(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("login")
		return ERR_ARGS, nil
	}

	epoch := time.Now().Unix()

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			login, confFile, passwd, description string
			nbf, exp                             int64
			enable, passStdin, help              bool
		)

		addFs := flag.NewFlagSet("add", flag.ExitOnError)
		addFs.StringVar(&login, "login", "", "username for login")
		addFs.BoolVar(&passStdin, "P", false, "stdin input password")
		addFs.StringVar(&passwd, "pass", "", "input password")

		addFs.Int64Var(&nbf, "nbf", epoch, "not before")
		addFs.Int64Var(&exp, "exp", 0, "expiration")
		addFs.BoolVar(&enable, "enable", false, "enable")
		addFs.StringVar(&description, "description", "", "description")
		addFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		addFs.BoolVar(&help, "help", false, "help me !!!")

		if err := addFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("login_add")
			return OK, nil
		}

		if len(addFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", addFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if nbf < epoch {
			return ERR_ARGS, fmt.Errorf("nbf must be greater than or equal to %d", epoch)
		}

		if exp != 0 && exp <= nbf {
			return ERR_ARGS, fmt.Errorf("exp must be greater than or equal to nbf -> %d", epoch)
		}

		if passStdin && passwd != "" {
			return ERR_ARGS, fmt.Errorf("error: -P or -pass are required")
		}

		login = strings.TrimSpace(login)

		// Vérification des args obligatoires
		if login == "" || (!passStdin && passwd == "") {
			addFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login and -P or -pass are required at least")
		}

		if err := pkgLogin.IsValidLogin(login); err != nil {
			return ERR_ARGS, err
		}

		if passStdin {
			var err error
			passwd, err = password()
			if err != nil {
				return ERR_ARGS, err
			}
		}

		// Générer le mot de passe
		pwdByte, err := argon2id.Hash(strings.TrimSpace(passwd), argon2id.DefaultParams())
		if err != nil {
			return ERR_PASSWORD, err
		}

		if err := pkgLogin.AddLogin(login, string(pwdByte), strings.TrimSpace(description), nbf, exp, !enable, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "del":

		var (
			login, confFile, uuid string
			help                  bool
		)

		delFs := flag.NewFlagSet("add", flag.ExitOnError)
		delFs.StringVar(&login, "login", "", "username for login")
		delFs.StringVar(&uuid, "uuid", "", "UUID for login")
		delFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		delFs.BoolVar(&help, "help", false, "help me !!!")

		if err := delFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("login_del")
			return OK, nil
		}

		if len(delFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", delFs.Args()[0])
		}

		// Vérification des args
		if (login == "" && uuid == "") || (login != "" && uuid != "") {
			delFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login OR -uuid must defined")
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if uuid != "" {
			return ERR_ARGS, pkgLogin.DelLoginByUUID(uuid, *cfg)
		}

		login = strings.TrimSpace(login)
		if err := pkgLogin.IsValidLogin(login); err != nil {
			return ERR_ARGS, err
		}

		if err := pkgLogin.DelLogin(login, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "get":

		var (
			login, confFile, uuid string
			all, help             bool
		)

		getFs := flag.NewFlagSet("get", flag.ExitOnError)
		getFs.StringVar(&login, "login", "", "username for login")
		getFs.StringVar(&uuid, "uuid", "", "uuid for login")
		getFs.BoolVar(&all, "all", false, "get all logins")
		getFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		getFs.BoolVar(&help, "help", false, "help me !!!")

		if err := getFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("login_get")
			return OK, nil
		}

		if len(getFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", getFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if all {
			if err := pkgLogin.GetLogins(*cfg); err != nil {
				return ERR_LOGIN, err
			}
			return OK, nil
		}

		if login != "" {
			if err := pkgLogin.GetLogin(login, *cfg); err != nil {
				return ERR_LOGIN, err
			}
			return OK, nil
		}

		if uuid != "" {
			if err := pkgLogin.GetLoginByUUID(uuid, *cfg); err != nil {
				return ERR_LOGIN, err
			}
			return OK, nil
		}

		getFs.Usage()
		return ERR_ARGS, fmt.Errorf("error: -all  or -login or -uuid must defined")

	case "check": // Vérifier l'authentification

		var (
			login, passwd, confFile string
			passStdin, help         bool
		)

		checkFs := flag.NewFlagSet("add", flag.ExitOnError)
		checkFs.StringVar(&login, "login", "", "username for login")
		checkFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		checkFs.BoolVar(&passStdin, "P", false, "stdin input password")
		checkFs.StringVar(&passwd, "pass", "", "input password")
		checkFs.BoolVar(&help, "help", false, "help me !!!")

		if err := checkFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("login_get")
			return OK, nil
		}

		if len(checkFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", checkFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if passStdin && passwd != "" {
			return ERR_ARGS, fmt.Errorf("error: -P or -pass are required")
		}

		login = strings.TrimSpace(login)

		// Vérification des args obligatoires
		if login == "" || (!passStdin && passwd == "") {
			checkFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login and -P or -pass are required at least")
		}

		if err := pkgLogin.IsValidLogin(login); err != nil {
			return ERR_LOGIN, err
		}

		if passStdin {
			var ctrl bool
			// Saisir le mot de passe
			fmt.Print("Password : ")
			passwd, ctrl, err = readPwdWithMask()
			if err != nil {
				return ERR_ARGS, fmt.Errorf("input password error : %s", err.Error())
			}
			fmt.Println("")
			if ctrl {
				return ERR_ARGS, fmt.Errorf("input password canceled")
			}
		}

		if err := pkgLogin.Check(login, passwd, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "enable", "disable":

		var (
			login, confFile, uuid string
			help                  bool
		)

		statusFs := flag.NewFlagSet("status", flag.ExitOnError)
		statusFs.StringVar(&login, "login", "", "username for login")
		statusFs.StringVar(&uuid, "uuid", "", "uuid for login")
		statusFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		statusFs.BoolVar(&help, "help", false, "help me !!!")

		if err := statusFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("login_add")
			return OK, nil
		}

		if len(statusFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", statusFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if !(login != uuid) {
			statusFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login or -uuid must defined")
		}

		if login != "" {
			var errLogin error = nil
			if arg == "enable" {
				errLogin = pkgLogin.EnableLogin(login, *cfg)
			} else {
				errLogin = pkgLogin.DisableLogin(login, *cfg)
			}
			if errLogin != nil {
				return ERR_LOGIN, errLogin
			}
		}

		if uuid != "" {
			var errUUID error = nil
			if arg == "enable" {
				errUUID = pkgLogin.EnableLoginByUUID(uuid, *cfg)
			} else {
				errUUID = pkgLogin.DisableLoginByUUID(uuid, *cfg)
			}
			if errUUID != nil {
				return ERR_LOGIN, errUUID
			}
		}

	default:
		fmt.Fprintf(os.Stdout, "command %s not exits\n", arg)

	}

	return OK, nil
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
