package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
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
