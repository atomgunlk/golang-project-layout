package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	// "github.com/jackc/pgx/v5"
)

//go:generate mockery --name Repository
type Repository interface {
	// Close
	Close() error
}

type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string

	OperationTimeout int // second
}

type repository struct {
	// DB               *pgx.Conn
	operationTimeout time.Duration
}

func New(cfg *Config) (Repository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port,
	)

	// conn, err := pgx.Connect(context.Background(), dsn)
	// if err != nil {
	// 	return nil, errors.Join(err, errors.New("[repository.New]"))
	// }

	// err = conn.Ping(context.Background())
	// if err != nil {
	// 	return nil, errors.New("[repository.New]: can not connect to DB")
	// }

	// Auto migration
	err = autoMigrator(cfg)
	if err != nil {
		return nil, err
	}

	operationTimeout := time.Duration(cfg.OperationTimeout) * time.Second
	if operationTimeout == 0 {
		return nil, errors.New("[repository.New]: operation timeout is invalid")
	}

	repo := &repository{DB: conn, operationTimeout: operationTimeout}
	return repo, nil
}

func (r *repository) defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), r.operationTimeout)
}

func (r *repository) Close() error {
	return r.DB.Close(context.Background())
}
