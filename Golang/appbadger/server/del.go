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

func DelServer(name, uuid, ip string, config pkgConfig.Config) error {

	pkgDefine.TrimSpace(&name, &uuid, &ip)

	if pkgDefine.CountDefinedStrArg(name, ip, uuid) != 1 {
		return fmt.Errorf("name or uuid or ip must defined")
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

	args := db.DelServerParams{
		Name: name,
		IP:   ipType,
		UUID: pgtype.UUID{Valid: false},
	}

	isSyntaxUUID := false
	pgUUID, err := pkgDefine.StringToPgUUID(uuid)
	if err == nil {
		args.Name = ""
		args.IP = netip.Addr{}
		args.UUID = pgUUID
		isSyntaxUUID = true
	}

	dbx := db.New(cnx.Cnx)
	rslt, err := dbx.DelServer(cnx.Ctx, args)
	if err != nil {

		if isSyntaxUUID {
			return json.NewEncoder(os.Stdout).Encode(struct {
				UUID   string
				Status bool
				Err    string
			}{
				UUID:   uuid,
				Status: false,
				Err:    err.Error(),
			})
		}

		if ipType.IsValid() {
			return json.NewEncoder(os.Stdout).Encode(struct {
				IP     netip.Addr
				Status bool
				Err    string
			}{
				IP:     ipType,
				Status: false,
				Err:    err.Error(),
			})
		}

		return json.NewEncoder(os.Stdout).Encode(struct {
			Name   string
			Status bool
			Err    string
		}{
			Name:   name,
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
		Name:   rslt.Name,
		IP:     rslt.Ip,
		Status: true,
	})

}
