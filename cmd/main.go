package main

import (
	"os"

	"github.com/SaidovZohid/personal-blog-task/api"
	"github.com/SaidovZohid/personal-blog-task/config"
	"github.com/SaidovZohid/personal-blog-task/database"
	"github.com/SaidovZohid/personal-blog-task/pkg/logger"
	"github.com/SaidovZohid/personal-blog-task/storage"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	_ "github.com/SaidovZohid/personal-blog-task/api/docs"
)

func main() {
	cfg := config.New()

	logger := logger.New(cfg.Environment)

	err := cfg.Load(".")
	if err != nil {
		logger.Panic("cannot load config", zap.Error(err))
	}

	// try to connect and migrate DB using golang-migrate.
	pgURL := cfg.PgURL()
	if err = database.MigrateDB(pgURL); err != nil {
		logger.Error("cannot migrate db", zap.Error(err))
		os.Exit(1)
	}

	// this returns connection pool
	dbPool, err := sqlx.Connect("postgres", pgURL)
	if err != nil {
		logger.Error("Unable to connect to database ", zap.Error(err))
		os.Exit(1)
	}
	defer dbPool.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	strg := storage.NewStorage(dbPool)
	inMemory := storage.NewInMemoryStorage(rdb)

	router := api.New(&api.RouterOptions{
		Cfg:      cfg,
		Storage:  strg,
		InMemory: inMemory,
		Logger:   logger,
	})

	if err := router.Run(cfg.RestAddr); err != nil {
		logger.Error("Unable to run rest api ", zap.Error(err))
	}
}
