package domain

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	"encoding/json"
	"os"
)

func AddDomain(domain string, config pkgConfig.Config) error {

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
	_, err = dbx.AddDomain(cnx.Ctx, dom)
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
