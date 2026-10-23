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
	pkgUser "appbadger/user"

	"github.com/matyas-cyril/argon2id"
	"golang.org/x/term"
)

func user(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("user")
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

		addFs := flag.NewFlagSet(arg, flag.ExitOnError)
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
			printHelp("user_add")
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

		if err := pkgUser.IsValidLogin(login); err != nil {
			return ERR_ARGS, err
		}

		if passStdin {
			var err error
			passwd, err = InputPassword()
			if err != nil {
				return ERR_ARGS, err
			}
		}

		// Générer le mot de passe
		pwdByte, err := argon2id.Hash(strings.TrimSpace(passwd), argon2id.DefaultParams())
		if err != nil {
			return ERR_PASSWORD, err
		}

		if err := pkgUser.AddLogin(login, string(pwdByte), strings.TrimSpace(description), nbf, exp, enable, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "del":

		var (
			login, confFile, uuid string
			help                  bool
		)

		delFs := flag.NewFlagSet(arg, flag.ExitOnError)
		delFs.StringVar(&login, "login", "", "username for login")
		delFs.StringVar(&uuid, "uuid", "", "UUID for login")
		delFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		delFs.BoolVar(&help, "help", false, "help me !!!")

		if err := delFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_del")
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
			return ERR_ARGS, pkgUser.DelLoginByUUID(uuid, *cfg)
		}

		login = strings.TrimSpace(login)
		if err := pkgUser.IsValidLogin(login); err != nil {
			return ERR_ARGS, err
		}

		if err := pkgUser.DelLogin(login, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "get":

		var (
			login, confFile, uuid string
			inactive, all, help   bool
		)

		getFs := flag.NewFlagSet(arg, flag.ExitOnError)
		getFs.StringVar(&login, "login", "", "username for login")
		getFs.StringVar(&uuid, "uuid", "", "uuid for login")
		getFs.BoolVar(&all, "all", false, "get all logins")
		getFs.BoolVar(&inactive, "inactive", false, "filter only disabled logins")
		getFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		getFs.BoolVar(&help, "help", false, "help me !!!")

		if err := getFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_get")
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

			if inactive {
				if err := pkgUser.GetLoginsInactive(*cfg); err != nil {
					return ERR_LOGIN, err
				}

			} else {
				if err := pkgUser.GetLogins(*cfg); err != nil {
					return ERR_LOGIN, err
				}

			}

			return OK, nil
		}

		if login != "" {

			if inactive {
				if err := pkgUser.GetLoginInactive(login, *cfg); err != nil {
					return ERR_LOGIN, err
				}

			} else {
				if err := pkgUser.GetLogin(login, *cfg); err != nil {
					return ERR_LOGIN, err
				}

			}

			return OK, nil
		}

		if uuid != "" {
			if err := pkgUser.GetLoginByUUID(uuid, *cfg); err != nil {
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

		checkFs := flag.NewFlagSet(arg, flag.ExitOnError)
		checkFs.StringVar(&login, "login", "", "username for login")
		checkFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		checkFs.BoolVar(&passStdin, "P", false, "stdin input password")
		checkFs.StringVar(&passwd, "pass", "", "input password")
		checkFs.BoolVar(&help, "help", false, "help me !!!")

		if err := checkFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_get")
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

		if err := pkgUser.IsValidLogin(login); err != nil {
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

		if err := pkgUser.Check(login, passwd, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "enable", "disable":

		var (
			login, confFile, uuid string
			help                  bool
		)

		statusFs := flag.NewFlagSet(arg, flag.ExitOnError)
		statusFs.StringVar(&login, "login", "", "username for login")
		statusFs.StringVar(&uuid, "uuid", "", "uuid for login")
		statusFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		statusFs.BoolVar(&help, "help", false, "help me !!!")

		if err := statusFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_add")
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
				errLogin = pkgUser.EnableLogin(login, *cfg)
			} else {
				errLogin = pkgUser.DisableLogin(login, *cfg)
			}
			if errLogin != nil {
				return ERR_LOGIN, errLogin
			}
		}

		if uuid != "" {
			var errUUID error = nil
			if arg == "enable" {
				errUUID = pkgUser.EnableLoginByUUID(uuid, *cfg)
			} else {
				errUUID = pkgUser.DisableLoginByUUID(uuid, *cfg)
			}
			if errUUID != nil {
				return ERR_LOGIN, errUUID
			}
		}

	case "modify":

		var (
			login, uuid, passwd, confFile, description string
			nbf, exp                                   int64
			passStdin, help, enable, disable           bool
		)

		modFs := flag.NewFlagSet(arg, flag.ExitOnError)
		modFs.StringVar(&login, "login", "", "username for login")
		modFs.StringVar(&uuid, "uuid", "", "uuid for login")
		modFs.BoolVar(&passStdin, "P", false, "stdin input password")
		modFs.StringVar(&passwd, "pass", "", "input password")
		modFs.Int64Var(&nbf, "nbf", 0, "not before")
		modFs.Int64Var(&exp, "exp", 0, "expiration")
		modFs.StringVar(&description, "description", "", "description")
		modFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		modFs.BoolVar(&enable, "enable", false, "enable user account")
		modFs.BoolVar(&disable, "disable", false, "disable user account")
		modFs.BoolVar(&help, "help", false, "help me !!!")

		if err := modFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_modify")
			return OK, nil
		}

		if len(modFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", modFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		uuid = strings.TrimSpace(uuid)
		login = strings.TrimSpace(login)

		if (uuid != "" && login != "") || (uuid == "" && login == "") {
			modFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login or -uuid are required")
		}

		if login != "" { // Login défini
			if err := pkgUser.IsValidLogin(login); err != nil {
				return ERR_ARGS, err
			}
		}

		if nbf == 0 && exp == 0 && passwd == "" && !passStdin && description == "" && !enable && !disable {
			modFs.Usage()
			return ERR_ARGS, fmt.Errorf("option required")
		}

		if enable && disable { // Ou tu sors ou j'te sors mais va falloir prendre une décision !
			return ERR_ARGS, fmt.Errorf("error: -enable OR -disable")
		}

		var disabled int8 = 0
		if enable {
			disabled = -1
		}

		if disable {
			disabled = 1
		}

		if nbf != 0 && nbf < epoch {
			return ERR_ARGS, fmt.Errorf("nbf must be greater than or equal to %d", epoch)
		}

		if exp != 0 && exp <= epoch {
			return ERR_ARGS, fmt.Errorf("exp must be greater than or equal to nbf -> %d", epoch)
		}

		if passStdin && passwd != "" {
			return ERR_ARGS, fmt.Errorf("error: -P or -pass are required")
		}

		if passStdin {
			var err error
			passwd, err = InputPassword()
			if err != nil {
				return ERR_ARGS, err
			}
		}

		var PassHash []byte
		if passwd != "" {
			// Générer le mot de passe
			PassHash, err = argon2id.Hash(strings.TrimSpace(passwd), argon2id.DefaultParams())
			if err != nil {
				return ERR_PASSWORD, err
			}
		}

		if err := pkgUser.ModifyLogin(login, uuid, string(PassHash), strings.TrimSpace(description), nbf, exp, disabled, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	case "expiration":
		var (
			login, uuid, confFile  string
			exp                    int64
			now, utc, remove, help bool
		)

		expFs := flag.NewFlagSet(arg, flag.ExitOnError)
		expFs.StringVar(&login, "login", "", "username for login")
		expFs.StringVar(&uuid, "uuid", "", "uuid for login")
		expFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		expFs.BoolVar(&now, "now", false, fmt.Sprintf("add epoch now(%d) value", time.Now().Unix()))
		expFs.Int64Var(&exp, "exp", 0, "add own epoch value")
		expFs.BoolVar(&remove, "remove", false, "remove expiratiion value")
		expFs.BoolVar(&utc, "utc", false, "heure UTC")
		expFs.BoolVar(&help, "help", false, "help me !!!")

		if err := expFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_expiration")
			return OK, nil
		}

		if len(expFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", expFs.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		uuid = strings.TrimSpace(uuid)
		login = strings.TrimSpace(login)

		if (uuid != "" && login != "") || (uuid == "" && login == "") {
			expFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -login or -uuid are required")
		}

		if login != "" { // Login défini
			if err := pkgUser.IsValidLogin(login); err != nil {
				return ERR_ARGS, err
			}
		}

		if (now && remove && exp != 0) || (!now && !remove && exp == 0) || (now && exp != 0) { // Faut faire eun choix
			return ERR_ARGS, fmt.Errorf("error: -now OR -remove OR -exp")
		}

		var srcEpoch int64

		if utc {
			srcEpoch = time.Now().UTC().Unix()
		} else {
			srcEpoch = time.Now().Unix()
		}

		var epoch int64 = 0
		if now {
			epoch = srcEpoch
		}

		if exp != 0 {
			if exp < srcEpoch {
				return ERR_ARGS, fmt.Errorf("error: -exp value %d is lower than %d", exp, srcEpoch)
			}
			epoch = exp
		}

		if err := pkgUser.ExpLogin(login, uuid, epoch, remove, *cfg); err != nil {
			return ERR_LOGIN, err
		}

		return OK, nil

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exits\n", arg)
		return ERR_ARGS, nil

	}

	return OK, nil
}

func InputPassword() (string, error) {

	// Traitement pour saisir le mot de passe
	fmt.Print("Password : ")
	pwd, ctrl, err := readPwdWithMask()
	if err != nil {
		return "", fmt.Errorf("input password error : %s", err.Error())
	}
	fmt.Println("")
	if ctrl {
		return "", fmt.Errorf("input password canceled")
	}

	// Vérifier la taille min et max du mot de passe
	pwd = strings.TrimSpace(pwd)
	if err := pkgUser.IsValidPassword(pwd); err != nil {
		return "", err
	}

	fmt.Print("Confirm  : ")
	confirmPwd, ctrl, err := readPwdWithMask()
	if err != nil {
		return "", fmt.Errorf("confirm password error : %s", err.Error())
	}
	fmt.Println("")
	if ctrl {
		return "", fmt.Errorf("input confirm password canceled")
	}

	if strings.Compare(pwd, confirmPwd) != 0 {
		return "", fmt.Errorf("input and confirm password mismatch")
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
		if pkgUser.IsValidPasswordRune(r) {
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
