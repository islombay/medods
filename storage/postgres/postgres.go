package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"medods/config"
	"medods/pkg/logs"
	"medods/storage"
	"os"
	"time"
)

type Storage struct {
	log logs.LoggerInterface
	db  *pgxpool.Pool
	cfg *config.ConfigDB

	auth    storage.AuthI
	user    storage.UserI
	session storage.SessionI
}

func New(cfg config.ConfigDB, log logs.LoggerInterface) storage.StorageInterface {
	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s sslmode=%s",
		cfg.Host,
		os.Getenv("DB_USER"),
		cfg.DBName,
		os.Getenv("DB_PWD"),
		cfg.Port,
		cfg.SSLMode,
	))
	if err != nil {
		panic(err)
	}

	conf.MaxConns = 40

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	return &Storage{
		db:  db,
		log: log,
		cfg: &cfg,

		auth:    NewAuthDB(db, log),
		user:    NewUserDB(db, log),
		session: NewSessionDB(db, log),
	}
}

func (db *Storage) Auth() storage.AuthI {
	return db.auth
}

func (db *Storage) User() storage.UserI {
	return db.user
}

func (db *Storage) Session() storage.SessionI {
	return db.session
}

func (db *Storage) Close() {

}

func (db *Storage) Migrate() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PWD"),
		db.cfg.Host, db.cfg.Port, db.cfg.DBName, db.cfg.SSLMode,
	)

	migrationsPath := fmt.Sprintf("file://%s", db.cfg.MigrationsPath)

	db.log.Debug("initializing migrations")

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Panic("could not initialize migration", logs.Error(err))
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Panic("could not migrate up", logs.Error(err))
		os.Exit(1)
	}

	db.log.Info("database migrated up")
}
