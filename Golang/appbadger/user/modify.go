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

func ModifyLogin(login, uuid, password, description string, nbf, exp int64, config pkgConfig.Config) error {

	if (uuid != "" && login != "") || (uuid == "" && login == "") {
		return fmt.Errorf("login or uuid must defined to modify")

	}
	var (
		pgUUID pgtype.UUID
		err    error
	)

	if uuid != "" {
		pgUUID, err = StringToPgUUID(uuid)
		if err != nil {
			return err
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

	// Définir le comportement par défaut
	args := db.ModifyLoginParams{
		Password:    pgtype.Text{Valid: false},
		NotBefore:   pgtype.Timestamptz{Valid: false},
		Expiration:  pgtype.Timestamptz{Valid: false},
		Description: pgtype.Text{Valid: false},
		Login:       login,
		UUID:        pgUUID,
	}

	if exp > 0 {
		time.Unix(exp, 0)
		args.Expiration = pgtype.Timestamptz{
			Time:  time.Unix(exp, 0),
			Valid: true,
		}
	}

	if nbf > 0 {
		time.Unix(nbf, 0)
		args.NotBefore = pgtype.Timestamptz{
			Time:  time.Unix(nbf, 0),
			Valid: true,
		}
	}

	if description != "" {
		args.Description = pgtype.Text{
			String: description,
			Valid:  true,
		}
	}

	if password != "" {
		args.Password = pgtype.Text{
			String: password,
			Valid:  true,
		}
	}

	fmt.Println(description)

	rslt, err := dbx.ModifyLogin(cnx.Ctx, args)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)
}
