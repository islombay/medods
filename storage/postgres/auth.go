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

func (db *authDB) GetByID(ctx context.Context, id string) (model.User, error) {
	q := `select
				id, first_name, last_name, hash_token,
				created_at, updated_at, deleted_at, email
			from users
			where id = $1 and deleted_at is null;`

	var user model.User

	if err := db.db.QueryRow(ctx, q, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.HashedRefreshToken,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Email,
	); err != nil {
		// TODO: handle error
		// Not found
		if errors.Is(err, pgx.ErrNoRows) {
			return user, error_list.NotFound
		}
		db.log.Error("could not find the user by id", logs.Error(err))
		return user, err
	}

	return user, nil
}

func (db *authDB) UpdateHash(ctx context.Context, user *model.User) error {
	q := `update users
			set hash_token = $1, updated_at = now()
			where id = $2 and deleted_at is null;`

	if _, err := db.db.Exec(ctx, q, user.HashedRefreshToken, user.ID); err != nil {
		db.log.Error("could not update hash token for user", logs.Error(err), logs.String("user_id", user.ID))
		return err
	}
	return nil
}
