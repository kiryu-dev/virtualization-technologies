package pgrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"virtualization-technologies/conf"
)

func NewPgConnection(cfg conf.DbConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.String())
	if err != nil {
		return nil, errors.WithMessage(err, "pgx connect")
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, errors.WithMessage(err, "postgres ping")
	}
	return conn, nil
}
