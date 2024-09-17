package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tender-service/config"
	"tender-service/internal/handler"
	"tender-service/internal/repository"
	"tender-service/internal/repository/postgres"
	"tender-service/internal/server"
	"tender-service/internal/service"
)

func Run(configPath string) {

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		logrus.Fatalf("Config error: %s", err)
	}

	SetLogrus(cfg.Log.Level)

	logrus.Info("Initializing postgres...")
	db, err := postgres.NewPostgresDB(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer db.Close()

	logrus.Info("Initializing repositories...")
	repos := repository.NewRepository(db)

	logrus.Info("Initializing services...")
	services := service.NewService(repos)

	logrus.Info("Initializing handlers and routes...")
	handlers := handler.NewHandler(services)

	logrus.Info("Starting http server...")
	logrus.Debugf("Server address: %s", cfg.HTTP.Address)
	srv := server.NewServer(handlers.InitRoutes(), server.Address(cfg.HTTP.Address))

	logrus.Info("Configuring graceful shutdown...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-quit:
		logrus.Info("app - Run - signal: " + s.String())
	case err = <-srv.Notify():
		logrus.Error(fmt.Errorf("app - Run - server.Notify: %w", err))
	}

	logrus.Info("Shutting down...")
	err = srv.Shutdown()
	if err != nil {
		logrus.Error(fmt.Errorf("app - Run - server.Shutdown: %w", err))
	}
}
