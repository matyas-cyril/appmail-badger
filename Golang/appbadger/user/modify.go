package user

import (
	pkgConfig "appbadger/config"
	db "appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ModifyLogin(login, uuid, password, description string, nbf, exp int64, disable int8, config pkgConfig.Config) error {

	pkgDefine.TrimSpace(&login, &uuid)

	if (uuid != "" && login != "") || (uuid == "" && login == "") {
		return fmt.Errorf("login or uuid must defined to modify")

	}
	var (
		pgUUID pgtype.UUID
		err    error
	)

	if uuid != "" {
		pgUUID, err = pkgDefine.StringToPgUUID(uuid)
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

	// Définir le comportement par défaut
	args := db.ModifyLoginParams{
		Password:    pgtype.Text{Valid: false},
		NotBefore:   pgtype.Timestamptz{Valid: false},
		Expiration:  pgtype.Timestamptz{Valid: false},
		Disable:     pgtype.Bool{Valid: false},
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
			String: strings.TrimSpace(description),
			Valid:  true,
		}
	}

	if password != "" {
		args.Password = pgtype.Text{
			String: password,
			Valid:  true,
		}
	}

	if disable != 0 {

		if disable == 1 {
			args.Disable = pgtype.Bool{
				Bool:  true,
				Valid: true,
			}
		} else {
			args.Disable = pgtype.Bool{
				Bool:  false,
				Valid: true,
			}
		}

	}

	rslt, err := dbx.ModifyLogin(cnx.Ctx, args)
	if err != nil {

		if login != "" {
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

		if uuid != "" {
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

	}

	rsltNbf := rslt.Nbf.Time.Unix()
	rsltExp := rslt.Exp.Time.Unix()

	if rsltNbf < 0 {
		rsltNbf = 0
	}

	if rsltExp < 0 {
		rsltExp = 0
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID        string
		Login       string
		Enable      bool
		Nbf         int64
		Exp         int64
		Description string
		Status      bool
	}{
		UUID:        rslt.Uuid.String(),
		Login:       rslt.Login,
		Enable:      !rslt.Disabled,
		Nbf:         rsltNbf,
		Exp:         rsltExp,
		Description: rslt.Description.String,
		Status:      true,
	})
}
