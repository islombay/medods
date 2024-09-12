package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"medods/internal/model"
	"medods/pkg/logs"
)

type userDB struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewUserDB(db *pgxpool.Pool, log logs.LoggerInterface) *userDB {
	return &userDB{
		db:  db,
		log: log,
	}
}

func (db *userDB) Create(ctx context.Context, user *model.User) error {
	q := `insert into users (first_name, last_name) values ($1, $2)
			returning id, created_at, updated_at, deleted_at`

	if err := db.db.QueryRow(ctx, q, user.FirstName, user.LastName).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	); err != nil {
		db.log.Error("Could not create user", logs.Error(err))
		return err
	}

	return nil
}
