package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"

	"github.com/matyas-cyril/argon2id"
)

func Check(login, password string, config pkgConfig.Config) error {

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

	rslt, err := dbx.GetPassword(cnx.Ctx, login)
	if err != nil {
		return err
	}

	check, err := argon2id.Check(password, rslt)
	if err != nil {
		return err
	}

	if check {
		return json.NewEncoder(os.Stdout).Encode(
			map[string]any{
				"login":  login,
				"check":  "authentication success",
				"status": 1,
			},
		)
	}

	return json.NewEncoder(os.Stdout).Encode(
		map[string]any{
			"login":  login,
			"check":  "authentication failure",
			"status": 0,
		},
	)

}
