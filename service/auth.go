package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"medods/internal/error_list"
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
	accessToken, refreshToken, err := jwt.GeneratePairTokens(m.UserId, m.IP)
	if err != nil {
		srv.log.Error("could not generate pair tokens", logs.Error(err))
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
	// add user to database
	user := model.User{
		FirstName: &m.FirstName,
		LastName:  &m.LastName,
		Email:     &m.Email,
	}
	if err := srv.storage.User().Create(ctx, &user); err != nil {
		return nil, err
	}

	// return token (call previous function)
	return srv.Login(ctx, model.LoginRequest{
		UserId: user.ID,
		IP:     m.IP,
	})
}

func (srv *Auth) Refresh(ctx context.Context, m model.RefreshRequest) (*model.TokenPair, error) {
	// parse access token
	tokenClaims, err := jwt.ParseToken(m.AccessToken)
	if err != nil {
		srv.log.Debug("could not parse jwt token", logs.Error(err))
		return nil, error_list.Unauthorized
	}

	// check user exists in db
	user, err := srv.storage.Auth().GetByID(ctx, tokenClaims.UserID)
	if err != nil {
		return nil, err
	}

	// check refresh token hash with db
	if err := bcrypt.CompareHashAndPassword([]byte(*user.HashedRefreshToken), []byte(m.RefreshToken)); err != nil {
		srv.log.Debug("trying to refresh token with invalid refresh token")
		return nil, error_list.Unauthorized
	}

	// TODO: check current IP address with IP address from jwt token if incorrect send email warning
	if m.IP != tokenClaims.IP {
		// TODO: send email warning
	}

	// generate token pairs
	accessToken, refreshToken, err := jwt.GeneratePairTokens(user.ID, m.IP)
	if err != nil {
		srv.log.Error("could not generate pair tokens", logs.Error(err))
		return nil, err
	}

	// save new hash to user
	hashedRefreshToken, err := helper.HashPassword(refreshToken)
	if err != nil {
		srv.log.Error("could not hash refresh token", logs.Error(err))
		return nil, err
	}

	user.HashedRefreshToken = &hashedRefreshToken
	if err := srv.storage.Auth().UpdateHash(ctx, &user); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
