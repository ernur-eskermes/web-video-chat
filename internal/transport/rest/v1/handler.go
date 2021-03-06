package v1

import (
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/olahol/melody.v1"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		h.initUsersRoutes(v1)
		h.initChatsRoutes(v1)
		h.initRoomsRoutes(v1)
	}
}

func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Should be min than " + fe.Param()
	case "max":
		return "Should be max than " + fe.Param()
	}

	return "Unknown error"
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
