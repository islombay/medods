package storage

import (
	"context"
	"medods/internal/model"
)

type StorageInterface interface {
	Close()
	Migrate()

	User() UserI
	Session() SessionI
}

type UserI interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type SessionI interface {
	Create(ctx context.Context, session *model.Session) error
	UpdateHash(ctx context.Context, session *model.Session) error
	GetByID(ctx context.Context, id string) (*model.Session, error)
}
