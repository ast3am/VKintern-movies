package config

import (
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/ast3am/VKintern-movies/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

var cfg *models.Config

func GetConfig(path string) *models.Config {
	logger := logging.GetLogger("")
	logger.DebugMsg("read config")
	cfg = &models.Config{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		logger.FatalMsg("", err)
	}
	return cfg
}
