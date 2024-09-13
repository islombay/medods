package storage

import (
	"context"
	"medods/internal/model"
)

type StorageInterface interface {
	Close()
	Migrate()

	Auth() AuthI
	User() UserI
}

type AuthI interface {
	UpdateHash(ctx context.Context, user *model.User) error
}

type UserI interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
}
