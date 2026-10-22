package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ExpLogin(login, uuid string, epoch int64, remove bool, config pkgConfig.Config) error {

	if (uuid != "" && login != "") || (uuid == "" && login == "") {
		return fmt.Errorf("login or uuid must defined to expiration")

	}
	var (
		pgUUID pgtype.UUID
		err    error
	)

	if uuid != "" {
		pgUUID, err = StringToPgUUID(uuid)
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

	var rslt any

	if epoch > 0 {

		args := db.SetExpirationParams{
			Expiration: pgtype.Timestamptz{Valid: false},
			Login:      login,
			UUID:       pgUUID,
		}

		time.Unix(epoch, 0)
		args.Expiration = pgtype.Timestamptz{
			Time:  time.Unix(epoch, 0),
			Valid: true,
		}

		rslt, err = dbx.SetExpiration(cnx.Ctx, args)
		if err != nil {
			return err
		}

	}

	if remove {

		args := db.ClearExpirationParams{
			Login: login,
			UUID:  pgUUID,
		}

		rslt, err = dbx.ClearExpiration(cnx.Ctx, args)
		if err != nil {
			return err
		}

	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}
