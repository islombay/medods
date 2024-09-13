package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"medods/pkg/logs"
)

type authDB struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewAuthDB(db *pgxpool.Pool, log logs.LoggerInterface) *authDB {
	return &authDB{
		db:  db,
		log: log,
	}
}
