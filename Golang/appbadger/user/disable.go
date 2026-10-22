package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
)

func DisableLogin[T string | pgtype.UUID](user T, config pkgConfig.Config) error {

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
		rslt, err := dbx.DisableLoginByUUID(cnx.Ctx, value)
		if err != nil {
			return json.NewEncoder(os.Stdout).Encode(struct {
				UUID   string
				Status bool
				Error  string
			}{
				UUID:   value.String(),
				Status: false,
				Error:  err.Error(),
			})
		}

		return json.NewEncoder(os.Stdout).Encode(struct {
			UUID    string
			Login   string
			Disable bool
			Status  bool
		}{
			UUID:    rslt.Uuid.String(),
			Login:   rslt.Login,
			Disable: rslt.Disabled,
			Status:  true,
		})

	}

	// On a forcément un string
	rslt, err := dbx.DisableLogin(cnx.Ctx, any(user).(string))
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Error  string
			Status bool
		}{
			Login:  any(user).(string),
			Error:  err.Error(),
			Status: false,
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID    string
		Login   string
		Disable bool
		Status  bool
	}{
		UUID:    rslt.Uuid.String(),
		Login:   rslt.Login,
		Disable: rslt.Disabled,
		Status:  true,
	})

}

func DisableLoginByUUID(uuid string, config pkgConfig.Config) error {
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
	return DisableLogin(pgUUID, config)
}

func EnableLogin[T string | pgtype.UUID](user T, config pkgConfig.Config) error {

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
		rslt, err := dbx.EnableLoginByUUID(cnx.Ctx, value)
		if err != nil {
			return json.NewEncoder(os.Stdout).Encode(struct {
				UUID   string
				Status bool
				Error  string
			}{
				UUID:   value.String(),
				Status: false,
				Error:  err.Error(),
			})
		}

		return json.NewEncoder(os.Stdout).Encode(struct {
			UUID   string
			Login  string
			Enable bool
			Status bool
		}{
			UUID:   rslt.Uuid.String(),
			Login:  rslt.Login,
			Enable: !rslt.Disabled,
			Status: true,
		})

	}

	// On a forcément un string
	rslt, err := dbx.EnableLogin(cnx.Ctx, any(user).(string))
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Error  string
			Status bool
		}{
			Login:  any(user).(string),
			Error:  err.Error(),
			Status: false,
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Login  string
		Enable bool
		Status bool
	}{
		UUID:   rslt.Uuid.String(),
		Login:  rslt.Login,
		Enable: !rslt.Disabled,
		Status: true,
	})

}

func EnableLoginByUUID(uuid string, config pkgConfig.Config) error {
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
	return EnableLogin(pgUUID, config)
}
