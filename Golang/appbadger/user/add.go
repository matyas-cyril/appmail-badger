package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func AddLogin(login, password, description string, nbf, exp int64, disable bool, config pkgConfig.Config) error {

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
		NotBefore:   pgtype.Timestamptz{Valid: false},
		Expiration:  pgtype.Timestamptz{Valid: false},
		Description: pgtype.Text{Valid: false},
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

	args.Username = login
	args.Password = password
	args.Disabled = disable

	rslt, err := dbx.AddLogin(cnx.Ctx, args)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}
