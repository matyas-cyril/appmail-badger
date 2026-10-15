package main

import (
	pkgDate "appbadger/date"
	"flag"
	"fmt"
	"os"
	"strings"
)

func date(args []string) (int, error) {

	if len(args) == 0 {
		printHelp("date")
		return ERR_ARGS, nil
	}

	switch arg := strings.TrimSpace(args[0]); arg {

	case "epoch":

		var (
			value     int64
			utc, help bool
		)
		epochFs := flag.NewFlagSet("epoch", flag.ExitOnError)
		epochFs.Int64Var(&value, "value", 0, "heure au format epoch")
		epochFs.BoolVar(&utc, "utc", false, "heure UTC")
		epochFs.BoolVar(&help, "help", false, "help me !!!")

		if err := epochFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("date_epoch")
			return OK, nil
		}

		if len(epochFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", epochFs.Args()[0])
		}

		if err := pkgDate.Epoch(value, utc); err != nil {
			return ERR_EPOCH, err
		}

		return OK, nil

	case "conv":

		var (
			day, hour string
			utc, help bool
		)

		convFs := flag.NewFlagSet("conv", flag.ExitOnError)
		convFs.StringVar(&day, "day", "", "jour au format YYYY/MM/DD")
		convFs.StringVar(&hour, "hour", "00:00:00", "jour au format HH:MM:SS")
		convFs.BoolVar(&utc, "utc", false, "heure UTC")
		convFs.BoolVar(&help, "help", false, "help me !!!")

		if err := convFs.Parse(args[1:]); err != nil {
			return ERR_ARGS, err
		}

		if help {
			printHelp("date_conv")
			return OK, nil
		}

		if len(convFs.Args()) > 0 {
			return ERR_ARGS, fmt.Errorf("unexpected arg : %q", convFs.Args()[0])
		}

		if err := pkgDate.Conv(day, hour, utc); err != nil {
			return ERR_EPOCH, err
		}

		return OK, nil

	default:
		fmt.Fprintf(os.Stdout, "command %s not exits\n", arg)

	}

	return OK, nil
}
