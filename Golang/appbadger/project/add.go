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

/*
	type ADD struct {
		Name string
		Sender string
		AuthUsers []string
		Domains []string
		ReplyTo string
		RcptMax uint16
		Burst uint16
		RateLimit string
		SizeMax uint64
		RcptPattern string
		Enable bool

}
*/
type ADD struct {
	Name      string
	Sender    string
	AuthUsers []string
	Domains   []string
	Enable    bool
}

func AddProject(args ADD, config pkgConfig.Config) error {

	name := strings.TrimSpace(args.Name)
	sender := strings.TrimSpace(args.Sender)
	authUsers := pkgDefine.TrimSpaceArray(args.AuthUsers)
	domains := pkgDefine.TrimSpaceArray(args.Domains)

	if name == "" || sender == "" || len(authUsers) == 0 || len(domains) == 0 {
		return fmt.Errorf("name AND sender AND auth user AND domains must defined")
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

	tx, err := cnx.Cnx.Begin(cnx.Ctx)
	if err != nil {
		return fmt.Errorf("failed begin transaction: %w", err)
	}
	defer tx.Rollback(cnx.Ctx)

	dbx := db.New(cnx.Cnx).WithTx(tx)

	argsAddProj := db.AddProjectParams{
		Name: name,
		Sender: pgtype.Text{
			String: sender,
			Valid:  true,
		},
		Enable: args.Enable,
	}

	// Liste des requêtes
	rsltAddProj, err := dbx.AddProject(cnx.Ctx, argsAddProj)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
	}

	argsAddDomProj := db.AddDomainInProjectParams{
		UUID:    rsltAddProj.Uuid,
		Domains: domains,
	}

	rsltAddDomProj, err := dbx.AddDomainInProject(cnx.Ctx, argsAddDomProj)
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
	}

	if len(rsltAddDomProj) == 0 {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    fmt.Errorf("dom %s not available", domains[:]).Error(),
		})
	}

	// Commit
	if err := tx.Commit(cnx.Ctx); err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
	}

	return json.NewEncoder(os.Stdout).Encode(struct {
		UUID   string
		Status bool
	}{
		UUID:   rsltAddProj.Uuid.String(),
		Status: true,
	})
}

/*
SELECT P.id_project AS uuid, P.name AS name, U.sender AS sender, P.reply_to AS reply_to, P.rcpt_max AS rcpt_max, P.burst AS burts, P.rate_limit AS rate_limit, P.size_max AS size_max, P.rcpt_pattern AS rcpt_pattern, P.modified AS modified, P.enabled AS enabled
FROM Projects P
INNER JOIN Senders U ON U.id_sender = P.id_sender



*/
