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

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"gopkg.in/olahol/melody.v1"

	"github.com/ernur-eskermes/web-video-chat/internal/config"
	"github.com/ernur-eskermes/web-video-chat/internal/repository"
	"github.com/ernur-eskermes/web-video-chat/internal/server"
	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/ernur-eskermes/web-video-chat/internal/transport/rest"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/ernur-eskermes/web-video-chat/pkg/database/mongodb"
	"github.com/ernur-eskermes/web-video-chat/pkg/hash"
	"github.com/ernur-eskermes/web-video-chat/pkg/logger"
	"github.com/ernur-eskermes/web-video-chat/pkg/room"
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
	OAuthInit(cfg)

	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Error(err)

		return
	}

	websocket := melody.New()

	db := mongoClient.Database(cfg.Mongo.Name)

	hasher := hash.NewSHA256Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	roomService := room.NewRoom(cfg.LiveKit.Host, cfg.LiveKit.ApiKey, cfg.LiveKit.ApiSecret)

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
		Room:            roomService,
		Websocket:       websocket,
	})
	handlers := rest.NewHandler(services, tokenManager, websocket)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running rest server: %s\n", err.Error())
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

func OAuthInit(cfg *config.Config) {
	store := sessions.NewCookieStore([]byte(cfg.Auth.SessionSecret))
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = cfg.Environment == config.Prod

	gothic.Store = store

	goth.UseProviders(
		google.New(cfg.GoogleOauth.ClientId, cfg.GoogleOauth.ClientSecret, cfg.GoogleOauth.CallbackURL),
	)
}
