package http

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"
	"net/http"

	"github.com/ernur-eskermes/web-video-chat/docs"
	"github.com/ernur-eskermes/web-video-chat/internal/config"
	v1 "github.com/ernur-eskermes/web-video-chat/internal/delivery/http/v1"
	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/ernur-eskermes/web-video-chat/pkg/limiter"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	websocket    *melody.Melody
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, websocket *melody.Melody) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
		websocket:    websocket,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)
	router.LoadHTMLGlob("templates/*")

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}

	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager, h.websocket)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
