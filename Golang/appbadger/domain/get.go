package domain

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	"encoding/json"
	"os"
)

func GetDomain(domain string, config pkgConfig.Config) error {

	dom, err := validDomain(domain)
	if err != nil {
		return err
	}

	cnx, err := db.NewPG(db.PGServer{
		Host:     config.Host,
		Database: config.DBName,
		User:     config.User,
		Password: config.Password,
		Port:     config.Port,
	})
	if err != nil {
		return err
	}
	defer cnx.ClosePG()

	dbx := db.New(cnx.Cnx)
	_, err = dbx.GetDomain(cnx.Ctx, dom)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Dom    string
			Status bool
			Err    string
		}{
			Dom:    dom,
			Status: false,
			Err:    err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Dom    string
		Status bool
	}{
		Dom:    dom,
		Status: true,
	})
}

func GetDomains(config pkgConfig.Config) error {

	cnx, err := db.NewPG(db.PGServer{
		Host:     config.Host,
		Database: config.DBName,
		User:     config.User,
		Password: config.Password,
		Port:     config.Port,
	})
	if err != nil {
		return err
	}
	defer cnx.ClosePG()

	dbx := db.New(cnx.Cnx)
	rslt, err := dbx.GetDomains(cnx.Ctx)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Dom    []string
		Status bool
	}{
		Dom:    rslt,
		Status: true,
	})
}
