package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Psql struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
}

func NewPsqlDb(user, password, host, port, db string) *Psql {
	return &Psql{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Db:       db,
	}
}

func (p *Psql) Connect() (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.User, p.Password, p.Host, p.Port, p.Db)
	return pgxpool.New(context.Background(), connStr)
}
