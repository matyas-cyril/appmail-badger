package domain

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	"encoding/json"
	"os"
)

func ModifyDomain(old, new string, config pkgConfig.Config) error {

	domOld, err := validDomain(old)
	if err != nil {
		return err
	}

	domNew, err := validDomain(new)
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
	_, err = dbx.ModifyDomain(cnx.Ctx, db.ModifyDomainParams{
		New: domNew,
		Old: domOld,
	})
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Old    string
			New    string
			Status bool
			Err    string
		}{
			Old:    domOld,
			New:    domNew,
			Status: false,
			Err:    err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Old    string
		New    string
		Status bool
	}{
		Old:    domOld,
		New:    domNew,
		Status: true,
	})
}
