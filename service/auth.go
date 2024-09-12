package service

import (
	"context"
	"medods/internal/model"
	"medods/pkg/helper"
	"medods/pkg/jwt"
	"medods/pkg/logs"
	"medods/storage"
)

type Auth struct {
	storage storage.StorageInterface
	log     logs.LoggerInterface
}

func NewAuth(storage storage.StorageInterface, log logs.LoggerInterface) *Auth {
	return &Auth{
		storage: storage,
		log:     log,
	}
}

func (srv *Auth) Login(ctx context.Context, m model.LoginRequest) (*model.TokenPair, error) {
	// check if user exists in database, else return unauthorized error
	user, err := srv.storage.Auth().GetByID(ctx, m.UserId)
	if err != nil {
		// Errors from storage layer will be logged by storage layer
		return nil, err
	}

	srv.log.Debug("user exists in database", logs.String("user_id", user.ID))

	// TODO: generate access and refresh tokens
	accessToken, err := jwt.Generate(user.ID, m.IP)
	if err != nil {
		srv.log.Error("Could not generate access token", logs.Error(err))
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		srv.log.Error("Could not generate refresh token", logs.Error(err))
		return nil, err
	}

	hashedRefreshToken, err := helper.HashPassword(refreshToken)
	if err != nil {
		srv.log.Error("Could not hash refresh token", logs.Error(err))
		return nil, err
	}

	user.HashedRefreshToken = &hashedRefreshToken

	// TODO: save hashed token to db
	if err := srv.storage.Auth().UpdateHash(ctx, &user); err != nil {
		return nil, err
	}

	srv.log.Debug("new hash saved for user database")

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *Auth) Register(ctx context.Context, m model.Register) (*model.TokenPair, error) {
	// TODO: add user to database
	user := model.User{
		FirstName: &m.FirstName,
		LastName:  &m.LastName,
	}
	if err := srv.storage.User().Create(ctx, &user); err != nil {
		return nil, err
	}

	// TODO: return token (call previous function)
	return srv.Login(ctx, model.LoginRequest{
		UserId: user.ID,
		IP:     m.IP,
	})
}
