package main

import (
	"github.com/gin-gonic/gin"
	"medods/api"
	"medods/config"
	"medods/pkg/logs"
	"medods/service"
	"medods/storage/postgres"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	log := logs.NewLogger("medods-auth", logs.LevelDebug)
	defer func() {
		if err := logs.Cleanup(log); err != nil {
			return
		}
	}()

	storage := postgres.New(cfg.DB, log)
	defer storage.Close()

	storage.Migrate()

	services := service.New(storage, log, cfg)

	r := gin.Default()

	api.New(r, services, cfg, log)

	go func() {
		if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
			log.Error("error listening host", logs.Error(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
