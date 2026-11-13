package project

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type GET struct {
	UUID    string
	Name    string
	Sender  string
	ReplyTo string
}

func GetProject(input GET, config pkgConfig.Config) (err error) {

	defer func() {
		if pErr := recover(); pErr != nil {
			err = json.NewEncoder(os.Stdout).Encode(struct {
				Status bool
				Err    string
			}{
				Status: false,
				Err:    fmt.Errorf("panic error : %s", pErr).Error(),
			})
		}
	}()

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

	args := db.GetProjectParams{
		UUID: pgtype.UUID{
			Valid: false,
		},
		Name:   pgtype.Empty.String(),
		Sender: pgtype.Empty.String(),
		ReplyTo: pgtype.Text{
			Valid: false,
		},
	}

	var cpt uint8 = 0

	pgUUID, err := pkgDefine.StringToPgUUID(strings.TrimSpace(input.UUID))
	if err == nil {
		args.UUID = pgUUID
		cpt++
	}

	if strings.TrimSpace(input.Name) != "" {
		args.Name = strings.TrimSpace(input.Name)
		cpt++
	}

	if strings.TrimSpace(input.Sender) != "" {
		args.Sender = strings.TrimSpace(input.Sender)
		cpt++
	}

	if strings.TrimSpace(input.ReplyTo) != "" {
		args.Sender = strings.TrimSpace(input.ReplyTo)
		args.ReplyTo = pgtype.Text{
			String: strings.TrimSpace(input.ReplyTo),
			Valid:  true,
		}
		cpt++
	}

	if cpt == 0 || cpt > 1 {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    fmt.Errorf("uuid or name or sender or reply_to or all must defined").Error(),
		})
	}

	rslt, err := dbx.GetProject(cnx.Ctx, args)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Error  string
		}{
			Status: false,
			Error:  err.Error(),
		})
	}

	type P struct {
		UUID        string
		Name        string
		Sender      string
		ReplyTo     string
		RcptMax     uint16
		Burst       uint16
		SizeMax     uint64
		RcptPattern string
		Enable      bool
		IPs         any
		NameServers any
		Domains     any
	}

	projects := []P{}

	for _, r := range rslt {

		p := P{
			UUID:        r.Uuid.String(),
			Name:        r.Name,
			Sender:      r.Sender,
			ReplyTo:     r.ReplyTo.String,
			RcptMax:     uint16(r.RcptMax.Int16),
			Burst:       uint16(r.Burst),
			SizeMax:     r.SizeMax.Int.Uint64(),
			Enable:      r.Enable,
			IPs:         r.Ips,
			NameServers: r.NameServers,
			Domains:     r.Domains,
		}

		projects = append(projects, p)

	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Projects []P
		Count    int
		Status   bool
	}{
		Projects: projects,
		Count:    len(projects),
		Status:   true,
	})

}

func GetProjects(config pkgConfig.Config) error {

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

	rslt, err := dbx.GetProjects(cnx.Ctx)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Error  string
		}{
			Status: false,
			Error:  err.Error(),
		})
	}

	type P struct {
		UUID        string
		Name        string
		Sender      string
		ReplyTo     string
		RcptMax     uint16
		Burst       uint16
		SizeMax     uint64
		RcptPattern string
		Enable      bool
		IPs         any
		NameServers any
		Domains     any
	}

	projects := []P{}

	for _, r := range rslt {

		p := P{
			UUID:        r.Uuid.String(),
			Name:        r.Name,
			Sender:      r.Sender,
			ReplyTo:     r.ReplyTo.String,
			RcptMax:     uint16(r.RcptMax.Int16),
			Burst:       uint16(r.Burst),
			SizeMax:     r.SizeMax.Int.Uint64(),
			Enable:      r.Enable,
			IPs:         r.Ips,
			NameServers: r.NameServers,
			Domains:     r.Domains,
		}

		projects = append(projects, p)

	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		Projects []P
		Count    int
		Status   bool
	}{
		Projects: projects,
		Count:    len(projects),
		Status:   true,
	})

}
