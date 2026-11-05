package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func DelLogin(id string, config pkgConfig.Config) error {

	id = strings.TrimSpace(id)
	if id == "" {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Error  string
		}{
			Status: false,
			Error:  "login or uuid not defined",
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

	args := db.DeleteLoginParams{
		Login: id,
		UUID:  pgtype.UUID{Valid: false},
	}

	isSyntaxUUID := false

	pgUUID, err := PkgDefine.StringToPgUUID(id)
	if err == nil {
		args.Login = ""
		args.UUID = pgUUID
		isSyntaxUUID = true
	}

	rslt, err := dbx.DeleteLogin(cnx.Ctx, args)
	if err != nil {

		if isSyntaxUUID {
			return json.NewEncoder(os.Stdout).Encode(struct {
				UUID   string
				Status bool
				Error  string
			}{
				UUID:   id,
				Status: false,
				Error:  err.Error(),
			})
		}

		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Status bool
			Error  string
		}{
			Login:  id,
			Status: false,
			Error:  err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Login  string
		Status bool
	}{
		UUID:   rslt.Uuid.String(),
		Login:  rslt.Login,
		Status: true,
	})

}
