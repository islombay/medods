package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"medods/internal/error_list"
	"medods/internal/mail"
	"medods/internal/model"
	"medods/pkg/helper"
	"medods/pkg/jwt"
	"medods/pkg/logs"
	"medods/storage"
)

type Auth struct {
	storage storage.StorageInterface
	log     logs.LoggerInterface
	mail    mail.MailServiceI
}

func NewAuth(storage storage.StorageInterface, log logs.LoggerInterface, mail mail.MailServiceI) *Auth {
	return &Auth{
		storage: storage,
		log:     log,
		mail:    mail,
	}
}

func (srv *Auth) Login(ctx context.Context, m model.LoginRequest) (*model.TokenPair, error) {
	// check if user exists in database, else return unauthorized error
	user, err := srv.storage.User().GetByID(ctx, m.UserId)
	if err != nil {
		// Errors from storage layer will be logged by storage layer
		return nil, err
	}

	srv.log.Debug("user exists in database", logs.String("user_id", user.ID))

	deviceInfo := ctx.Value("device_info").(model.DeviceInfo)

	// create new session
	session := model.Session{
		UserID: &user.ID,
		IP:     &deviceInfo.IP,
	}
	if err := srv.storage.Session().Create(ctx, &session); err != nil {
		return nil, err
	}

	// generate access and refresh tokens
	accessToken, refreshToken, err := jwt.GeneratePairTokens(m.UserId, deviceInfo.IP, session.ID)
	if err != nil {
		srv.log.Error("could not generate pair tokens", logs.Error(err))
		return nil, err
	}

	hashedRefreshToken, err := helper.HashPassword(refreshToken)
	if err != nil {
		srv.log.Error("Could not hash refresh token", logs.Error(err))
		return nil, err
	}

	session.Hash = &hashedRefreshToken

	// save hashed token to db
	if err := srv.storage.Session().UpdateHash(ctx, &session); err != nil {
		return nil, err
	}

	srv.log.Debug("new hash saved for session database")

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
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

	deviceInfo := ctx.Value("device_info").(model.DeviceInfo)

	// return token (call previous function)
	return srv.Login(ctx, model.LoginRequest{
		UserId: user.ID,
		IP:     deviceInfo.IP,
	})
}

func (srv *Auth) Refresh(ctx context.Context, m model.RefreshRequest) (*model.TokenPair, error) {
	// parse access token
	tokenClaims, err := jwt.ParseToken(m.AccessToken)
	if err != nil {
		if !errors.Is(err, error_list.TokenExpired) {
			srv.log.Debug("could not parse jwt token", logs.Error(err))
			return nil, error_list.Unauthorized
		}
	}
	// check user exists in db
	user, err := srv.storage.User().GetByID(ctx, tokenClaims.UserID)
	if err != nil {
		return nil, err
	}

	// get session by session_id
	session, err := srv.storage.Session().GetByID(ctx, tokenClaims.SessionID)
	if err != nil {
		return nil, err
	}

	// check refresh token hash with db
	if err := bcrypt.CompareHashAndPassword([]byte(*session.Hash), []byte(m.RefreshToken)); err != nil {
		srv.log.Debug("trying to refresh token with invalid refresh token")
		return nil, error_list.Unauthorized
	}

	deviceInfo := ctx.Value("device_info").(model.DeviceInfo)

	if deviceInfo.IP != tokenClaims.IP {
		// send email warning
		srv.log.Debug("refresh token requested from another IP",
			logs.String("got IP", deviceInfo.IP),
			logs.String("expected IP", tokenClaims.IP),
		)
		if user.Email != nil {
			go func() {
				if err := srv.mail.WarnIPAddressChange(*user.Email, deviceInfo); err != nil {
					srv.log.Error("could not send email notification about IP address change",
						logs.Error(err),
						logs.String("to", *user.Email),
					)
				}
			}()
		}
	}

	// generate token pairs
	accessToken, refreshToken, err := jwt.GeneratePairTokens(user.ID, deviceInfo.IP, session.ID)
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

	session.Hash = &hashedRefreshToken
	if err := srv.storage.Session().UpdateHash(ctx, session); err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
