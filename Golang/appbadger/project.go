package main

import (
	pkgConfig "appbadger/config"
	pkgDefine "appbadger/define"
	pkgProject "appbadger/project"
	"flag"
	"fmt"
	"strings"
)

func project(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("project")
		return ERR_ARGS, nil
	}

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			name, sender, domains, auths string
			confFile                     string
			enable, help                 bool
		)

		addProj := flag.NewFlagSet(arg, flag.ExitOnError)
		addProj.StringVar(&name, "name", "", "define name")
		addProj.StringVar(&sender, "sender", "", "define sender")
		addProj.StringVar(&domains, "domains", "", "define sender")
		addProj.StringVar(&auths, "auths", "", "define auth users")
		addProj.BoolVar(&enable, "enable", false, "enable project")
		addProj.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		addProj.BoolVar(&help, "help", false, "help me !!!")

		if err := addProj.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("user_add")
			return OK, nil
		}

		if len(addProj.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", addProj.Args()[0])
		}

		pkgDefine.TrimSpace(&name, &sender, &confFile, &domains, &auths)

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if err := pkgProject.AddProject(pkgProject.ADD{
			Name:      name,
			Sender:    sender,
			AuthUsers: strings.Split(auths, ","),
			Domains:   strings.Split(domains, ","),
			Enable:    enable,
		}, *cfg); err != nil {
			return ERR_SERVER, err
		}

		return OK, nil

	}

	return 0, nil

}
