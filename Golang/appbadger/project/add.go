package project

import (
	pkgConfig "appbadger/config"
	"appbadger/db"
	pkgDefine "appbadger/define"
	"encoding/json"
	"fmt"
	"math/big"
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
	Name        string
	Sender      string
	ReplyTo     string
	AuthUsers   []string
	Domains     []string
	Servers     []string
	RateLimit   string
	RcptMax     *uint16
	Burst       *uint16
	SizeMax     *uint64
	RcptPattern string // A définir le fonctionnement
	Enable      bool
}

func AddProject(args ADD, config pkgConfig.Config) error {

	name := strings.TrimSpace(args.Name)
	sender := strings.TrimSpace(args.Sender)
	replyTo := strings.TrimSpace(args.ReplyTo)
	authUsers := pkgDefine.TrimSpaceArray(args.AuthUsers)
	domains := pkgDefine.TrimSpaceArray(args.Domains)
	servers := parseServers(pkgDefine.TrimSpaceArray(args.Servers))
	rateLimit := strings.TrimSpace(args.RateLimit)
	rcptPattern := strings.TrimSpace(args.RcptPattern)

	if name == "" || sender == "" || len(authUsers) == 0 || len(domains) == 0 || servers == nil {
		return fmt.Errorf("name AND sender AND auth user AND domains AND servers must defined")
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
		Name:   name,
		Sender: sender,
		ReplyTo: pgtype.Text{
			String: sender,
			Valid:  true,
		},
		RateLimit: pgtype.Text{
			Valid: false,
		},
		RcptMax: pgtype.Int2{
			Valid: false,
		},
		Burst: pgtype.Int4{
			Valid: false,
		},
		SizeMax: pgtype.Numeric{
			Valid: false,
		},
		RcptPattern: pgtype.Text{
			Valid: false,
		},
		Enable: args.Enable,
	}

	if replyTo != "" {
		argsAddProj.ReplyTo = pgtype.Text{
			String: replyTo,
			Valid:  true,
		}
	}

	if rateLimit != "" {
		r, err := CheckRateLimitSyntax(rateLimit)
		if err != nil {
			return json.NewEncoder(os.Stdout).Encode(struct {
				Status bool
				Err    string
			}{
				Status: false,
				Err:    err.Error(),
			})
		}
		argsAddProj.RateLimit = pgtype.Text{
			String: r,
			Valid:  true,
		}
	}

	if args.RcptMax != nil {
		argsAddProj.RcptMax = pgtype.Int2{
			Int16: int16(*args.RcptMax),
			Valid: true,
		}
	}

	if args.Burst != nil {
		argsAddProj.Burst = pgtype.Int4{
			Int32: int32(*args.Burst),
			Valid: true,
		}
	}

	if args.SizeMax != nil {
		argsAddProj.SizeMax = pgtype.Numeric{
			Int:   new(big.Int).SetUint64(*args.SizeMax),
			Exp:   0, // entier
			Valid: true,
		}
	}

	if rcptPattern != "" {
		argsAddProj.RcptPattern = pgtype.Text{
			String: rcptPattern,
			Valid:  true,
		}
	}

	// Ajoute des informations dans le Project
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

	// Ajouter le projet à des domaines existants
	rsltAddDomProj, err := dbx.AddDomainInProject(cnx.Ctx, db.AddDomainInProjectParams{
		UUID:    rsltAddProj.Uuid,
		Domains: domains,
	})
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
			Err:    fmt.Errorf("at least one domain %s not available", domains[:]).Error(),
		})
	}

	// Ajouter des utilisateurs pour l'Authentification
	rsltAddAuthProj, err := dbx.AddAuthInProject(cnx.Ctx, db.AddAuthInProjectParams{
		UUID:  rsltAddProj.Uuid,
		Users: authUsers,
	})
	if err != nil {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    err.Error(),
		})
	}

	if len(rsltAddAuthProj) == 0 {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    fmt.Errorf("at least one auth user %s not available", authUsers[:]).Error(),
		})
	}

	// Ajouter des servers pour auth la connexion
	rsltAddSrvProj, err := dbx.AddServerInProject(cnx.Ctx, db.AddServerInProjectParams{
		UUIDProject: rsltAddProj.Uuid,
		IPs:         servers.ip,
		ServerNames: servers.name,
	})

	if len(rsltAddSrvProj) == 0 {
		return json.NewEncoder(os.Stdout).Encode(struct {
			Status bool
			Err    string
		}{
			Status: false,
			Err:    fmt.Errorf("at least one auth server %s not available", pkgDefine.TrimSpaceArray(args.Servers)[:]).Error(),
		})
	}

	// Commit si tout est OK, sinon Rollback
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
