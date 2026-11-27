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

// Définition identique
type MOD = ADD

func ModProject(uuid string, args MOD, config pkgConfig.Config) error {

	// Vérifier la syntaxe - UUID obligatoire car c'est la clef
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

	name := strings.TrimSpace(args.Name)
	sender := strings.TrimSpace(args.Sender)
	replyTo := strings.TrimSpace(args.ReplyTo)
	authUsers := pkgDefine.TrimSpaceArray(args.AuthUsers)
	domains := pkgDefine.TrimSpaceArray(args.Domains)
	servers := parseServers(pkgDefine.TrimSpaceArray(args.Servers))
	rateLimit := strings.TrimSpace(args.RateLimit)
	rcptPattern := strings.TrimSpace(args.RcptPattern)

	//	fmt.Println(name, sender, replyTo, authUsers, domains, servers, rateLimit, rcptPattern)
	fmt.Println(pgUUID)

	mod_proj := false
	mod_proj_dom := false
	mod_proj_srv := false

	argsModProj := db.ModProjectParams{
		Name: pgtype.Text{
			Valid: false,
		},
		Sender: pgtype.Text{
			Valid: false,
		},
		ReplyTo: pgtype.Text{
			Valid: false,
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
		Enable: pgtype.Bool{
			Valid: false,
		},
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

	if name != "" {
		argsModProj.Name = pgtype.Text{
			String: name,
			Valid:  true,
		}
		mod_proj = true
	}

	if sender != "" {
		argsModProj.Sender = pgtype.Text{
			String: sender,
			Valid:  true,
		}
		mod_proj = true
	}

	if replyTo != "" {
		argsModProj.ReplyTo = pgtype.Text{
			String: replyTo,
			Valid:  true,
		}
		mod_proj = true
	}

	if rcptPattern != "" {
		argsModProj.RcptPattern = pgtype.Text{
			String: rcptPattern,
			Valid:  true,
		}
		mod_proj = true
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
		argsModProj.RateLimit = pgtype.Text{
			String: r,
			Valid:  true,
		}
		mod_proj = true
	}

	if args.RcptMax != nil {
		argsModProj.RcptMax = pgtype.Int2{
			Int16: int16(*args.RcptMax),
			Valid: true,
		}
		mod_proj = true
	}

	if args.Burst != nil {
		argsModProj.Burst = pgtype.Int4{
			Int32: int32(*args.Burst),
			Valid: true,
		}
		mod_proj = true
	}

	if args.SizeMax != nil {
		argsModProj.SizeMax = pgtype.Numeric{
			Int:   new(big.Int).SetUint64(*args.SizeMax),
			Exp:   0, // entier
			Valid: true,
		}
		mod_proj = true
	}

	if mod_proj {
		_, err := dbx.ModProject(cnx.Ctx, argsModProj)
		if err != nil {
			return json.NewEncoder(os.Stdout).Encode(struct {
				Status bool
				Err    string
			}{
				Status: false,
				Err:    err.Error(),
			})
		}
	}

	if mod_proj_dom {
		// ProjectsDomains
	}

	if mod_proj_srv {
		// ProjectsServers
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
		UUID:   pgUUID.String(),
		Status: true,
	})
}
