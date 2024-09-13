package service

import (
	"medods/config"
	"medods/internal/mail"
	"medods/pkg/logs"
	"medods/storage"
)

type ServiceInterface interface {
	Auth() *Auth
	User() *User
}

type Service struct {
	cfg     config.Config
	log     logs.LoggerInterface
	storage storage.StorageInterface

	auth *Auth
	user *User
}

func New(storage storage.StorageInterface, log logs.LoggerInterface, cfg config.Config) ServiceInterface {
	srv := Service{
		cfg:     cfg,
		log:     log,
		storage: storage,
	}

	mailService := mail.NewMailService(log, cfg.Mail)

	srv.auth = NewAuth(storage, log, mailService)
	srv.user = NewUser(storage, log)

	return &srv
}

func (srv *Service) Auth() *Auth {
	return srv.auth
}

func (srv *Service) User() *User {
	return srv.user
}
