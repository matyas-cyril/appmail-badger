package main

import (
	pkgConfig "appbadger/config"
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
			printHelp("domain_add")
			return OK, nil
		}

		if len(addSrv.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", addSrv.Args()[0])
		}

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		name = strings.TrimSpace(name)
		ip = strings.TrimSpace(ip)
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

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exits\n", arg)
		return ERR_ARGS, nil

	}

	return OK, nil
}
