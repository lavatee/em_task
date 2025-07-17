package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lavatee/subs"
	"github.com/lavatee/subs/internal/endpoint"
	"github.com/lavatee/subs/internal/repository"
	"github.com/lavatee/subs/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Subscription Service API
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн-подписках пользователей
// @host localhost:8080
// @BasePath /api/v1

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	if err := InitConfig(); err != nil {
		logger.Fatalf("Failed to init config: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Fatalf("Failed to open Postgres DB: %s", err.Error())
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatalf("Failed to create migrate driver: %s", err.Error())
	}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationsPath := "file://" + filepath.Join(dir, "../schema")
	migrations, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logger.Fatalf("Failed to create migrate instance: %s", err.Error())
	}
	if err = migrations.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("Migrations error: %s", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo, logger)
	endp := endpoint.NewEndpoint(services, logger)
	server := &subs.Server{}
	go func() {
		if err := server.Run(viper.GetString("port"), endp.InitRoutes()); err != nil {
			logger.Fatalf("Failed to run server: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := server.Shutdown(); err != nil {
		logger.Fatalf("Failed to shutdown server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logger.Fatalf("Failed to disconnect Postgres DB: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
