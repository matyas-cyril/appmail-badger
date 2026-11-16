package project

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"os"
	"strings"
)

func EnableProject(uuid string, config pkgConfig.Config) error {
	return activeProject(uuid, true, config)
}

func DisableProject(uuid string, config pkgConfig.Config) error {
	return activeProject(uuid, false, config)
}

func activeProject(uuid string, enable bool, config pkgConfig.Config) error {

	pgUUID, err := pkgDefine.StringToPgUUID(strings.TrimSpace(uuid))
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
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

	rslt, err := dbx.EnableProjet(cnx.Ctx, db.EnableProjetParams{
		UUID:   pgUUID,
		Enable: enable,
	})
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Error  string
		}{
			Status: false,
			Error:  err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Name   string
		Sender string
		Enable bool
		Status bool
	}{
		UUID:   rslt.Uuid.String(),
		Name:   rslt.Name,
		Sender: rslt.Sender,
		Enable: rslt.Enabled,
		Status: true,
	})

}
