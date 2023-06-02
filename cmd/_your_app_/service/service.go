package service

import (
	"time"

	"github.com/atomgunlk/YOUR-REPO-NAME/cmd/YOUR-APP/repository"
)

// Config represents service configuration
type Config struct {
	AppEnv    string
	AppPort   int
	LogLevel  string
	JwtSecret string
}

type Dependency struct {
	Repo *repository.Repository
}

type Memory struct {
	saveImageCacheExpire map[string]time.Time
}

type Service struct {
	config *Config
	deps   *Dependency
	mem    *Memory
}

func New(cfg *Config, dep *Dependency) (*Service, error) {
	s := Service{
		config: cfg,
		deps:   dep,
		mem:    new(Memory),
	}

	// init service mem
	s.mem.saveImageCacheExpire = make(map[string]time.Time)

	return &s, nil
}
