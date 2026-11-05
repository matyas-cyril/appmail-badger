package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"os"

	"github.com/matyas-cyril/argon2id"
)

func Check(login, password string, config pkgConfig.Config) error {

	pkgDefine.TrimSpace(&login, &password)

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
		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Status bool
			Error  string
		}{
			Login:  login,
			Status: false,
			Error:  err.Error(),
		})
	}

	check, err := argon2id.Check(password, rslt)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Status bool
			Error  string
		}{
			Login:  login,
			Status: false,
			Error:  err.Error(),
		})
	}

	if check {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Login  string
			Auth   string
			Status bool
		}{
			Login:  login,
			Auth:   "authentication success",
			Status: true,
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Login  string
		Auth   string
		Status bool
	}{
		Login:  login,
		Auth:   "authentication failed",
		Status: false,
	})

}
