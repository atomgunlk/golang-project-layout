package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/handler"
	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/repository"
	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/service"
	"github.com/atomgunlk/golang-common/pkg/env"
	"github.com/atomgunlk/golang-common/pkg/graceful"
	"github.com/atomgunlk/golang-common/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	f := fiber.New()

	appEnv := env.RequiredEnv("APP_ENV")
	appPort, err := strconv.Atoi(env.RequiredEnv("APP_PORT"))
	if err != nil {
		logger.WithError(err).Panic("[main]: Invalid APP_PORT")
	}
	logLevel := env.RequiredEnv("LOG_LEVEL")
	jwtSecret := env.RequiredEnv("JWT_SECRET")

	dbHost := env.RequiredEnv("DB_HOST")
	dbPort, err := strconv.Atoi(env.RequiredEnv("DB_PORT"))
	if err != nil {
		logger.WithError(err).Panic("[main]: Invalid DB_PORT")
	}
	dbName := env.RequiredEnv("DB_NAME")
	dbUsername := env.RequiredEnv("DB_USERNAME")
	dbPassword := env.RequiredEnv("DB_PASSWORD")

	repo, err := repository.New(
		&repository.Config{
			Host:     dbHost,
			Port:     dbPort,
			Database: dbName,

			Username: dbUsername,
			Password: dbPassword,

			OperationTimeout: 10,
		},
	)
	if err != nil {
		logger.WithError(err).Panic("[main]: Unable to New repository")
	}

	s, err := service.New(&service.Config{
		AppEnv:    appEnv,
		AppPort:   appPort,
		LogLevel:  logLevel,
		JwtSecret: jwtSecret,
	}, &service.Dependency{
		Repo: repo
	})

	h, err := handler.New(
		&handler.Config{
			JWTSecret: jwtSecret,
		},
		&handler.Dependency{
			Service: s,
		},
	)
	if err != nil {
		logger.WithError(err).Panic("[main]: Unable to New handler")
	}

	err = h.InitRoutes(f)
	if err != nil {
		logger.WithError(err).Panic("[main]: Unable to init fiber routes")
	}

	go func() {
		if err := f.Listen(fmt.Sprintf(":%d", appPort)); err != nil {
			logger.Fatal(err)
		}
	}()

	if err := graceful.ListenSignal(func() error {

		err := f.ShutdownWithTimeout(30 * time.Second)
		if err != nil {
			logger.WithError(err).Error("[main]: Shutdown Fiber fail")
		}
		err = repo.Close()
		if err != nil {
			logger.WithError(err).Error("[main]: Shutdown repository fail")
			return err
		}

		logger.Info("[main]: Graceful shutdown success")
		return nil

	}); err != nil {
		logger.WithError(err).Panic("[main]: unable gracefully shutdown the service")
	}
}
