package server

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func ModifyServer(uuid, name, ip, description string, config pkgConfig.Config) error {

	pkgDefine.TrimSpace(&uuid, &name, &ip, &description)

	if uuid == "" {
		return fmt.Errorf("uuid must defined")
	}

	if pkgDefine.CountDefinedStrArg(name, ip, description) == 0 {
		return fmt.Errorf("at least name or uuid or ip or description must defined")
	}

	var ipType netip.Addr
	// On vérifie que l'on a une IP, si ip existe
	if strings.TrimSpace(ip) != "" {
		var errIp error
		ipType, errIp = netip.ParseAddr(ip)
		if errIp != nil {
			return errIp
		}
	}

	pgUUID, err := pkgDefine.StringToPgUUID(uuid)
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

	args := db.ModifyServerParams{
		UUID:        pgUUID,
		Name:        name,
		IP:          nil,
		Description: pgtype.Text{Valid: false},
	}

	if description != "" {
		args.Description = pgtype.Text{
			String: strings.TrimSpace(description),
			Valid:  true,
		}
	}

	if ipType.IsValid() {
		args.IP = &ipType
	}

	dbx := db.New(cnx.Cnx)
	rslt, err := dbx.ModifyServer(cnx.Ctx, args)
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

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID        string
		Name        string
		IP          netip.Addr
		Description string
		Status      bool
	}{
		UUID:        rslt.Uuid.String(),
		Name:        rslt.Name,
		IP:          rslt.Ip,
		Description: rslt.Description.String,
		Status:      true,
	})
}
