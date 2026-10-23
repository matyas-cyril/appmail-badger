package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	"encoding/json"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetLoginsInactive(config pkgConfig.Config) error {

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

	rslt, err := queries.GetLoginsInactive(cnx.Ctx)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Error  string
		}{
			Status: false,
			Error:  err.Error(),
		})
	}

	type L struct {
		UUID   string
		Login  string
		Enable bool
		Nbf    int64
		Exp    int64
	}

	logins := []L{}

	for _, r := range rslt {

		var (
			rsltNbf int64
			rsltExp int64
			ok      bool
		)

		rsltNbf, ok = r.Nbf.(int64)
		if !ok {
			rsltNbf = 0
		}

		rsltExp, ok = r.Exp.(int64)
		if !ok {
			rsltExp = 0
		}

		l := L{
			UUID:   r.Uuid.String(),
			Login:  r.Login,
			Enable: !r.Disabled,
			Nbf:    rsltNbf,
			Exp:    rsltExp,
		}

		logins = append(logins, l)

	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Inactive []L
		Status   bool
	}{
		Inactive: logins,
		Status:   true,
	})

}

func GetLoginInactive(id string, config pkgConfig.Config) error {

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

	args := db.GetLoginInactiveParams{
		Login: id,
		UUID:  pgtype.UUID{Valid: false},
	}

	isSyntaxUUID := false

	pgUUID, err := StringToPgUUID(id)
	if err == nil {
		args.Login = ""
		args.UUID = pgUUID
		isSyntaxUUID = true
	}

	rslt, err := dbx.GetLoginInactive(cnx.Ctx, args)
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

	var (
		rsltNbf int64
		rsltExp int64
		ok      bool
	)

	rsltNbf, ok = rslt.Nbf.(int64)
	if !ok {
		rsltNbf = 0
	}

	rsltExp, ok = rslt.Exp.(int64)
	if !ok {
		rsltExp = 0
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Login  string
		Enable bool
		Nbf    int64
		Exp    int64
		Status bool
	}{
		UUID:   rslt.Uuid.String(),
		Login:  rslt.Login,
		Enable: !rslt.Disabled,
		Nbf:    int64(rsltNbf),
		Exp:    int64(rsltExp),
		Status: true,
	})

}
