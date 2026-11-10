package main

import (
	pkgConfig "appbadger/config"
	pkgDefine "appbadger/define"
	pkgProject "appbadger/project"
	"flag"
	"fmt"
	"os"
	"strconv"
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
			name, sender, domains, auths     string
			rcptMaxStr, burstStr, sizeMaxStr string
			rateLimit, servers, confFile     string
			enable, help                     bool
		)

		addProj := flag.NewFlagSet(arg, flag.ExitOnError)
		addProj.StringVar(&name, "name", "", "define name")
		addProj.StringVar(&sender, "sender", "", "define sender")
		addProj.StringVar(&domains, "domains", "", "define sender")
		addProj.StringVar(&auths, "auths", "", "define auth users")
		addProj.StringVar(&servers, "servers", "", "define server name or ip")
		addProj.BoolVar(&enable, "enable", false, "enable project")
		addProj.StringVar(&rcptMaxStr, "rcpt_max", "", "define rcpt max")
		addProj.StringVar(&burstStr, "burst", "", "define burst")
		addProj.StringVar(&sizeMaxStr, "size_max", "", "define size max")
		addProj.StringVar(&rateLimit, "rate_limit", "", "define rate limit")
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

		pkgDefine.TrimSpace(&name, &sender, &confFile, &domains, &auths, &rateLimit)

		// Charger la conf
		cfg, err := pkgConfig.Load(confFile)
		if err != nil {
			return ERR_CONF_FILE, err
		}

		var rcptMax *uint16 = nil
		if rcptMaxStr = strings.TrimSpace(rcptMaxStr); rcptMaxStr != "" {
			val, err := strconv.ParseUint(rcptMaxStr, 10, 16)
			if err != nil {
				return ERR_CONF_FILE, fmt.Errorf("arg rcpt_max invalid integer: %w", err)
			}
			if val > 1500 {
				return ERR_CONF_FILE, fmt.Errorf("arg rcpt_max value must be between 0 and 1500")
			}
			rcptMax = new(uint16(val))
		}

		var burst *uint16 = nil
		if burstStr = strings.TrimSpace(burstStr); burstStr != "" {
			val, err := strconv.ParseUint(burstStr, 10, 16)
			if err != nil {
				return ERR_CONF_FILE, fmt.Errorf("arg rcpt_max invalid integer: %w", err)
			}
			burst = new(uint16(val))
		}

		var sizeMax *uint64 = nil
		if burstStr = strings.TrimSpace(burstStr); burstStr != "" {
			val, err := strconv.ParseUint(burstStr, 10, 64)
			if err != nil {
				return ERR_CONF_FILE, fmt.Errorf("arg rcpt_max invalid integer: %w", err)
			}
			sizeMax = &val
		}

		if err := pkgProject.AddProject(pkgProject.ADD{
			Name:      name,
			Sender:    sender,
			AuthUsers: strings.Split(auths, ","),
			Domains:   strings.Split(domains, ","),
			Servers:   strings.Split(servers, ","),
			Enable:    enable,
			RcptMax:   rcptMax,
			Burst:     burst,
			SizeMax:   sizeMax,
		}, *cfg); err != nil {
			return ERR_SERVER, err
		}

	case "get":

	case "del":

	case "enable", "disable":

	case "modify":

	default:
		fmt.Fprintf(os.Stdout, "command '%s' not exits\n", arg)
		return ERR_ARGS, nil

	}

	return OK, nil

}
