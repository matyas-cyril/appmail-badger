package login

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
)

func DelLogin[T string | pgtype.UUID](user T, config pkgConfig.Config) error {

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

	// On test uniquement le type pgtype.UUID, car après uniquement string
	if value, ok := any(user).(pgtype.UUID); ok {
		rslt, err := dbx.DeleteLoginByUUID(cnx.Ctx, value)
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(rslt)
	}

	// On a forcément un string
	rslt, err := dbx.DeleteLogin(cnx.Ctx, any(user).(string))
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(rslt)

}

func DelLoginByUUID(uuid string, config pkgConfig.Config) error {
	pgUUID, err := StringToPgUUID(uuid)
	if err != nil {
		return err
	}
	return DelLogin(pgUUID, config)
}
