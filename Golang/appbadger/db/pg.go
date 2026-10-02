package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PGServer struct {
	Host     string
	Database string
	User     string
	Password string
	Port     uint16
	SslMode  bool
}

type PG struct {
	Cnx *pgx.Conn
	Ctx context.Context
}

func NewPG(conf PGServer) (*PG, error) {

	// Définir des valeurs par défaut
	if conf.Port == 0 {
		conf.Port = 5432
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	sslmode := "disable"
	if conf.SslMode {
		sslmode = "enable"
	}

	// Créer l'uri
	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		strings.TrimSpace(conf.User),
		strings.TrimSpace(conf.Password),
		strings.TrimSpace(conf.Host),
		conf.Port,
		strings.TrimSpace(conf.Database),
		sslmode,
	)

	// Établir la connexion
	ctx := context.Background()
	cnx, err := pgx.Connect(ctx, uri)
	if err != nil {
		return nil, err
	}

	if err := cnx.Ping(ctx); err != nil {
		cnx.Close(ctx)
		return nil, fmt.Errorf("failed to check connection (ping) : %w", err)
	}

	return &PG{
		Cnx: cnx,
		Ctx: ctx,
	}, nil

}

func (p *PG) ClosePG() error {

	if p.Cnx == nil {
		return nil
	}

	if err := p.Cnx.Close(p.Ctx); err != nil {
		return fmt.Errorf("failed to close Postgres connection : %w", err)
	}

	p.Cnx = nil
	return nil
}
