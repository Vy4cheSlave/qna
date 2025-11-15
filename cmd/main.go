package main

import (
	// internal
	"github.com/Vy4cheSlave/qna/internal/config"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/db"
	"github.com/Vy4cheSlave/qna/internal/infrastructure/rest"
	"github.com/Vy4cheSlave/qna/internal/logpack"
	"github.com/Vy4cheSlave/qna/internal/usecase"
	// external
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	// std
	"context"
	// "fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal("Ошибка загрузки env файла:", err)
	}

	// Загрузка конфигурарции из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := logpack.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}
	restAddr := strings.Join([]string{cfg.Rest.Host, cfg.Rest.Port}, ":")

	// Подключение к базе данных
	repo, err := db.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing repository"))
	}

	// Инициализация сервиса
	service := usecase.NewQNAManagerService(repo, repo)

	// Запуск HTTP-сервера в отдельной горутине
	app := rest.NewApp(logger, &restAddr, service)
	go func() {
		err := app.ServerInstance.Run()
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	signalFromChannel := <-signalChan

	logger.Info("Shutting down server...", slog.String("signal", signalFromChannel.String()))
	logger.Info("Shutting down gracefully...")
}
