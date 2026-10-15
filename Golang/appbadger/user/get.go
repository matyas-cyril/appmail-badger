package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"
)

func GetLogin(login string, config pkgConfig.Config) error {

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

	queries := db.New(cnx.Cnx)

	rslt, err := queries.GetLogin(cnx.Ctx, login)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}

func GetLoginByUUID(uuid string, config pkgConfig.Config) error {

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

	pgUUID, err := StringToPgUUID(uuid)
	if err != nil {
		return err
	}

	dbx := db.New(cnx.Cnx)

	rslt, err := dbx.GetLoginByUUID(cnx.Ctx, pgUUID)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}

func GetLogins(config pkgConfig.Config) error {

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

	queries := db.New(cnx.Cnx)

	rslt, err := queries.GetLogins(cnx.Ctx)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}
