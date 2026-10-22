package main

import (
	pkgConfig "appbadger/config"
	pkgDomain "appbadger/domain"
	"flag"
	"fmt"
	"os"
	"strings"
)

func domain(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("domain")
		return ERR_ARGS, nil
	}

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			domain, confFile string
			help             bool
		)

		addDom := flag.NewFlagSet(arg, flag.ExitOnError)
		addDom.StringVar(&domain, "dom", "", "define a new domain")
		addDom.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		addDom.BoolVar(&help, "help", false, "help me !!!")

		if err := addDom.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("domain_add")
			return OK, nil
		}

		if len(addDom.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", addDom.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		domain = strings.TrimSpace(domain)
		// Vérification des args obligatoires
		if domain == "" {
			addDom.Usage()
			return ERR_ARGS, fmt.Errorf("error: -domain is required")
		}

		if err := pkgDomain.AddDomain(domain, *cfg); err != nil {
			return ERR_LOGIN, err
		}

	case "del":

		var (
			domain, confFile string
			help             bool
		)

		delDom := flag.NewFlagSet(arg, flag.ExitOnError)
		delDom.StringVar(&domain, "dom", "", "define a new domain")
		delDom.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		delDom.BoolVar(&help, "help", false, "help me !!!")

		if err := delDom.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("domain_del")
			return OK, nil
		}

		if len(delDom.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", delDom.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		domain = strings.TrimSpace(domain)
		// Vérification des args obligatoires
		if domain == "" {
			delDom.Usage()
			return ERR_ARGS, fmt.Errorf("error: -domain is required")
		}

		if err := pkgDomain.DelDomain(domain, *cfg); err != nil {
			return ERR_DOM, err
		}

		return OK, nil

	case "get":

		var (
			domain, confFile string
			help, all        bool
		)

		getDom := flag.NewFlagSet(arg, flag.ExitOnError)
		getDom.StringVar(&domain, "dom", "", "define a new domain")
		getDom.BoolVar(&all, "all", false, "get all domains")
		getDom.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		getDom.BoolVar(&help, "help", false, "help me !!!")

		if err := getDom.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("domain_get")
			return OK, nil
		}

		if len(getDom.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", getDom.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if all {
			if err := pkgDomain.GetDomains(*cfg); err != nil {
				return ERR_DOM, err
			}
			return OK, nil
		}

		domain = strings.TrimSpace(domain)
		// Vérification des args obligatoires
		if domain == "" {
			getDom.Usage()
			return ERR_ARGS, fmt.Errorf("error: -domain is required")
		}

		if err := pkgDomain.GetDomain(domain, *cfg); err != nil {
			return ERR_DOM, err
		}

	case "modify":

		var (
			old, new, confFile string
			help               bool
		)

		modDom := flag.NewFlagSet(arg, flag.ExitOnError)
		modDom.StringVar(&new, "new", "", "input the new name domain")
		modDom.StringVar(&old, "old", "", "input the old name domain")
		modDom.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		modDom.BoolVar(&help, "help", false, "help me !!!")

		if err := modDom.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("domain_get")
			return OK, nil
		}

		if len(modDom.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", modDom.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		old = strings.TrimSpace(old)
		new = strings.TrimSpace(new)

		// Vérification des args obligatoires
		if old == "" || new == "" {
			modDom.Usage()
			return ERR_ARGS, fmt.Errorf("error: -old AND -new must defined")
		}

		if old == new {
			modDom.Usage()
			return ERR_ARGS, fmt.Errorf("error: -old AND -new value must different")
		}

		if err := pkgDomain.ModifyDomain(old, new, *cfg); err != nil {
			return ERR_DOM, err
		}

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exits\n", arg)
		return ERR_ARGS, nil

	}

	return OK, nil
}
