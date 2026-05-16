package github

import (
	"github.com/Chandra179/gosdk/logger"
)

type Dependencies struct {
	Config Config
	Logger logger.Logger
}

func NewDependencies(cfg Config) *Dependencies {
	log := logger.NewLogger(cfg.Logger.Level)
	return &Dependencies{
		Config: cfg,
		Logger: log,
	}
}
