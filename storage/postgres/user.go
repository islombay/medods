package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"medods/internal/error_list"
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
	q := `insert into users (first_name, last_name, email) values ($1, $2, $3)
			returning id, created_at, updated_at, deleted_at`

	if err := db.db.QueryRow(ctx, q, user.FirstName, user.LastName, user.Email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	); err != nil {
		db.log.Error("Could not create user", logs.Error(err))
		return err
	}

	return nil
}

func (db *userDB) GetByID(ctx context.Context, id string) (*model.User, error) {
	q := `select
				id, first_name, last_name,
				created_at, updated_at, deleted_at, email
			from users
			where id = $1 and deleted_at is null;`

	var user model.User

	if err := db.db.QueryRow(ctx, q, id).Scan(
		&user.ID, &user.FirstName, &user.LastName,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Email,
	); err != nil {
		// TODO: handle error
		// Not found
		if errors.Is(err, pgx.ErrNoRows) {
			return &user, error_list.NotFound
		}
		db.log.Error("could not find the user by id", logs.Error(err))
		return &user, err
	}

	return &user, nil
}
