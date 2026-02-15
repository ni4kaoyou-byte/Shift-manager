package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/config"
	"github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	setGinMode(cfg.AppEnv)
	router := server.NewRouter()

	log.Info().Str("port", cfg.Port).Str("app_env", cfg.AppEnv).Msg("starting api server")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func setGinMode(appEnv string) {
	switch appEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
