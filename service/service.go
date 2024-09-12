package service

import (
	"medods/config"
	"medods/pkg/logs"
	"medods/storage"
)

type ServiceInterface interface {
	Auth() *Auth
}

type Service struct {
	cfg     config.Config
	log     logs.LoggerInterface
	storage storage.StorageInterface

	auth *Auth
}

func New(storage storage.StorageInterface, log logs.LoggerInterface, cfg config.Config) ServiceInterface {
	srv := Service{
		cfg:     cfg,
		log:     log,
		storage: storage,
	}

	srv.auth = NewAuth(storage, log)

	return &srv
}

func (srv *Service) Auth() *Auth {
	return srv.auth
}
