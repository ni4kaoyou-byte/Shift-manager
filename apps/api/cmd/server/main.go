package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Info().Str("port", cfg.Port).Str("app_env", cfg.AppEnv).Msg("starting api server")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}
