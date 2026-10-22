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

	var (
		rsltNbf int64
		rsltExp int64
		ok      bool
	)

	rsltNbf, ok = rslt.Nbf.(int64)
	if !ok || rsltNbf == 0 {
		rsltNbf = -1
	}

	rsltExp, ok = rslt.Exp.(int64)
	if !ok || rsltExp == 0 {
		rsltExp = -1
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID        string
		Login       string
		Enable      bool
		Nbf         int64
		Exp         int64
		Description string
		Status      bool
	}{
		UUID:        rslt.Uuid.String(),
		Login:       rslt.Login,
		Enable:      !rslt.Disabled,
		Nbf:         int64(rsltNbf),
		Exp:         int64(rsltExp),
		Description: rslt.Description.String,
		Status:      true,
	})

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

	dbx := db.New(cnx.Cnx)

	rslt, err := dbx.GetLoginByUUID(cnx.Ctx, pgUUID)
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

	var (
		rsltNbf int64
		rsltExp int64
		ok      bool
	)

	rsltNbf, ok = rslt.Nbf.(int64)
	if !ok || rsltNbf == 0 {
		rsltNbf = -1
	}

	rsltExp, ok = rslt.Exp.(int64)
	if !ok || rsltExp == 0 {
		rsltExp = -1
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID        string
		Login       string
		Enable      bool
		Nbf         int64
		Exp         int64
		Description string
		Status      bool
	}{
		UUID:        rslt.Uuid.String(),
		Login:       rslt.Login,
		Enable:      !rslt.Disabled,
		Nbf:         int64(rsltNbf),
		Exp:         int64(rsltExp),
		Description: rslt.Description.String,
		Status:      true,
	})

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
