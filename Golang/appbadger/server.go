package main

import (
	pkgConfig "appbadger/config"
	pkgDefine "appbadger/define"
	pkgServer "appbadger/server"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

func server(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("server")
		return ERR_ARGS, nil
	}

	switch arg := strings.TrimSpace(args[0]); arg {

	case "add":

		var (
			name, ip, confFile, description string
			help                            bool
		)

		addSrv := flag.NewFlagSet(arg, flag.ExitOnError)
		addSrv.StringVar(&name, "name", "", "define name of server ")
		addSrv.StringVar(&ip, "ip", "", "define ip of server ")
		addSrv.StringVar(&description, "description", "", "description")
		addSrv.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		addSrv.BoolVar(&help, "help", false, "help me !!!")

		if err := addSrv.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("server_add")
			return OK, nil
		}

		if len(addSrv.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", addSrv.Args()[0])
		}

		pkgDefine.TrimSpace(&name, &ip, &confFile)

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if name == "" || ip == "" {
			addSrv.Usage()
			return ERR_ARGS, fmt.Errorf("error: -name AND -ip are required")
		}

		// On vérifie que l'on a une IP
		ipType, err := netip.ParseAddr(ip)
		if err != nil {
			return ERR_ARGS, err
		}

		if err := pkgServer.AddServer(name, ipType, description, *cfg); err != nil {
			return ERR_SERVER, err
		}

	case "del":

		var (
			name, ip, confFile, uuid string
			help                     bool
		)

		delFs := flag.NewFlagSet(arg, flag.ExitOnError)
		delFs.StringVar(&name, "name", "", "define name of server ")
		delFs.StringVar(&ip, "ip", "", "define ip of server ")
		delFs.StringVar(&uuid, "uuid", "", "defined UUID for server")
		delFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		delFs.BoolVar(&help, "help", false, "help me !!!")

		if err := delFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("server_del")
			return OK, nil
		}

		if len(delFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", delFs.Args()[0])
		}

		pkgDefine.TrimSpace(&name, &ip, &uuid)

		// Vérification des args
		if pkgDefine.CountDefinedStrArg(name, ip, uuid) != 1 {
			delFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -name OR -ip OR -uuid must defined")
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if err := pkgServer.DelServer(name, uuid, ip, *cfg); err != nil {
			return ERR_SERVER, err
		}

	case "get":

		var (
			name, ip, confFile, uuid string
			help, all                bool
		)

		getFs := flag.NewFlagSet(arg, flag.ExitOnError)
		getFs.StringVar(&name, "name", "", "define name of server ")
		getFs.StringVar(&ip, "ip", "", "define ip of server ")
		getFs.StringVar(&uuid, "uuid", "", "defined UUID for server")
		getFs.BoolVar(&all, "all", false, "get all servers")
		getFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		getFs.BoolVar(&help, "help", false, "help me !!!")

		if err := getFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("server_get")
			return OK, nil
		}

		if len(getFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", getFs.Args()[0])
		}

		pkgDefine.TrimSpace(&name, &ip, &uuid)

		// Vérification des args
		if !all {
			if pkgDefine.CountDefinedStrArg(name, ip, uuid) != 1 {
				getFs.Usage()
				return ERR_ARGS, fmt.Errorf("error: -name OR -ip OR -uuid must defined")
			}
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if all {
			if err := pkgServer.GetServers(*cfg); err != nil {
				return ERR_SERVER, err
			}
		} else {
			if err := pkgServer.GetServer(name, uuid, ip, *cfg); err != nil {
				return ERR_SERVER, err
			}
		}

	case "modify":

		var (
			uuid, confFile, name, ip, description string
			help                                  bool
		)

		modFs := flag.NewFlagSet(arg, flag.ExitOnError)
		modFs.StringVar(&uuid, "uuid", "", "modify UUID for server")
		modFs.StringVar(&name, "name", "", "modify name of server ")
		modFs.StringVar(&ip, "ip", "", "modify ip of server ")
		modFs.StringVar(&description, "description", "", "modify description")
		modFs.StringVar(&confFile, "conf", DEFAULT_CONFIG_FILE, "json configuration file")
		modFs.BoolVar(&help, "help", false, "help me !!!")

		if err := modFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("server_modify")
			return OK, nil
		}

		if len(modFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", modFs.Args()[0])
		}

		pkgDefine.TrimSpace(&uuid, &name, &ip, &description, &confFile)

		// Obligation d'avoir l'UUID
		if uuid == "" {
			modFs.Usage()
			return ERR_ARGS, fmt.Errorf("error: -uuid must defined")
		}

		if pkgDefine.CountDefinedStrArg(name, ip, description) == 0 {
			return ERR_ARGS, fmt.Errorf("error: at least one argument (-name, -ip, -decription) must be defined")
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		if err := pkgServer.ModifyServer(uuid, name, ip, description, *cfg); err != nil {
			return ERR_SERVER, err
		}

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exits\n", arg)
		return ERR_ARGS, nil

	}

	return OK, nil
}
