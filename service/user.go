package service

import (
	"context"
	"medods/internal/model"
	"medods/pkg/logs"
	"medods/storage"
)

type User struct {
	storage storage.StorageInterface
	log     logs.LoggerInterface
}

func NewUser(storage storage.StorageInterface, log logs.LoggerInterface) *User {
	return &User{
		storage: storage,
		log:     log,
	}
}

func (srv *User) GetUser(ctx context.Context, user_id string) (*model.User, error) {
	user, err := srv.storage.User().GetByID(ctx, user_id)
	if err != nil {
		return user, err
	}
	return user, nil
}
