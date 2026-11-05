package server

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	"encoding/json"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func AddServer(name string, ip netip.Addr, description string, config pkgConfig.Config) error {

	name = strings.TrimSpace(name)

	if name == "" || !ip.IsValid() {
		return fmt.Errorf("name or ip invalid")
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

	args := db.AddServerParams{
		Name:        name,
		IP:          ip,
		Description: pgtype.Text{Valid: false},
	}

	if description != "" {
		args.Description = pgtype.Text{
			String: strings.TrimSpace(description),
			Valid:  true,
		}
	}

	dbx := db.New(cnx.Cnx)
	rslt, err := dbx.AddServer(cnx.Ctx, args)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Name   string
			IP     netip.Addr
			Status bool
			Err    string
		}{
			Name:   name,
			IP:     ip,
			Status: false,
			Err:    err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Name   string
		IP     netip.Addr
		Status bool
	}{
		UUID:   rslt.Uuid.String(),
		Name:   name,
		IP:     ip,
		Status: true,
	})

}
