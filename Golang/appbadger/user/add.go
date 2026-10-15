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

func AddLogin(login, passHash, description string, nbf, exp int64, enable bool, config pkgConfig.Config) error {

	if login == "" || passHash == "" {
		return fmt.Errorf("login or password must defined")
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

	args := db.AddLoginParams{
		Login:       login,
		Password:    passHash,
		NotBefore:   pgtype.Timestamptz{Valid: false},
		Expiration:  pgtype.Timestamptz{Valid: false},
		Description: pgtype.Text{Valid: false},
		Disabled: pgtype.Bool{
			Bool:  !enable,
			Valid: true,
		}}

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

	rslt, err := dbx.AddLogin(cnx.Ctx, args)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}
