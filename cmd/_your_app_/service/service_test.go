package service_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/repository"
	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/service"
	"github.com/jlc-group/golang-common/pkg/env"
	"github.com/jlc-group/golang-common/pkg/logger"
)

var (
	s *service.Service
)

func setup() {
	// init config
	dbPort, err := strconv.Atoi(env.RequiredEnv("DB_PORT"))
	if err != nil {
		logger.WithError(err).Panic("[test.Service]: Invalid DB_PORT")
	}

	jwtSecret = env.RequiredEnv("JWT_SECRET")

	r, err := repository.New(&repository.Config{
		Host:     env.RequiredEnv("DB_HOST"),
		Port:     dbPort,
		Username: env.RequiredEnv("DB_USERNAME"),
		Password: env.RequiredEnv("DB_PASSWORD"),
		Database: env.RequiredEnv("DB_NAME"),

		OperationTimeout: 10,
	})
	if err != nil {
		logger.WithError(err).Panic("[test.Service]: unable to new repository")
	}

	// Setup service
	s, err = service.New(
		&service.Config{
			AppEnv:    env.RequiredEnv("APP_ENV"),
			AppPort:   env.RequiredEnv("APP_PORT"),
			LogLevel:  env.RequiredEnv("LOG_LEVEL"),
			JwtSecret: env.RequiredEnv("JWT_SECRET"),
		},
		&service.Dependency{
			Repo: r,
		},
	)
	if err != nil {
		logger.WithError(err).Panic("[test.Service]: unable to new service")
	}
}
func shutdown() {
	// close client
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
