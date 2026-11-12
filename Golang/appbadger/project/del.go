package project

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"os"
)

func DelProject(uuid string, config pkgConfig.Config) error {

	pgUUID, err := pkgDefine.StringToPgUUID(uuid)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			UUID   string
			Status bool
			Error  string
		}{
			UUID:   uuid,
			Status: false,
			Error:  err.Error(),
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

	rslt, err := dbx.DelProjet(cnx.Ctx, pgUUID)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			UUID   string
			Status bool
			Error  string
		}{
			UUID:   pgUUID.String(),
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
