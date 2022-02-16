//nolint: funlen
package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ernur-eskermes/web-video-chat/internal/config"
	delivery "github.com/ernur-eskermes/web-video-chat/internal/delivery/http"
	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"github.com/ernur-eskermes/web-video-chat/internal/server"
	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/ernur-eskermes/web-video-chat/pkg/database/mongodb"
	"github.com/ernur-eskermes/web-video-chat/pkg/hash"
	"github.com/ernur-eskermes/web-video-chat/pkg/logger"
)

// @title Web-Video-Chat API
// @version 1.0
// @description REST API for Web-Video-Chat App

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:           repos,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
		Environment:     cfg.Environment,
		Domain:          cfg.HTTP.Host,
	})
	handlers := delivery.NewHandler(services, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
