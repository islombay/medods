package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"medods/internal/error_list"
	"medods/internal/model"
	"medods/pkg/logs"
	"medods/storage"
)

type sessionDB struct {
	db  *pgxpool.Pool
	log logs.LoggerInterface
}

func NewSessionDB(db *pgxpool.Pool, log logs.LoggerInterface) storage.SessionI {
	return &sessionDB{
		db:  db,
		log: log,
	}
}

func (db *sessionDB) Create(ctx context.Context, session *model.Session) error {
	q := `insert into sessions(user_id, ip) values ($1, $2)
			returning id, created_at, updated_at`

	if err := db.db.QueryRow(ctx, q, session.UserID, session.IP).Scan(
		&session.ID, &session.CreatedAt, &session.UpdatedAt,
	); err != nil {
		db.log.Error("could not create session", logs.Error(err))
		return err
	}

	return nil
}

func (db *sessionDB) UpdateHash(ctx context.Context, session *model.Session) error {
	q := `update sessions
			set hash = $1, updated_at = now()
			where id = $2 and deleted_at is null;`

	if _, err := db.db.Exec(ctx, q, session.Hash, session.ID); err != nil {
		db.log.Error("could not update hash token for session",
			logs.Error(err),
			logs.String("user_id", *session.UserID),
			logs.String("session_id", session.ID),
		)
		return err
	}
	return nil
}

func (db *sessionDB) GetByID(ctx context.Context, id string) (*model.Session, error) {
	q := `select
			id, hash, user_id, ip,
			created_at, updated_at, deleted_at
		from sessions
		where id = $1 and deleted_at is null`

	var session model.Session
	if err := db.db.QueryRow(ctx, q, id).Scan(
		&session.ID, &session.Hash, &session.UserID, &session.IP,
		&session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, error_list.NotFound
		}
		db.log.Error("could not get session by id",
			logs.Error(err),
			logs.String("session_id", id),
		)
		return nil, err
	}

	return &session, nil
}
