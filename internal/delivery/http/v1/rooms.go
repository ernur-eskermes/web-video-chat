package v1

import (
	"net/http"

	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) initRoomsRoutes(api *gin.RouterGroup) {
	rooms := api.Group("/rooms")
	{
		authenticated := rooms.Group("/", h.userIdentity)
		{
			authenticated.POST("/create", h.roomCreate)
		}
	}
}

type roomCreateInput struct {
	Name string `json:"name" binding:"required,min=8,max=64"`
}

type roomCreateResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}

// @Summary Room Create
// @Tags rooms-create
// @Security UsersAuth
// @Description create room
// @ModuleID roomCreate
// @Accept  json
// @Produce  json
// @Param input body roomCreateInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /rooms/create [post]
func (h *Handler) roomCreate(c *gin.Context) {
	var inp roomCreateInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	roomId, token, err := h.services.Rooms.Create(c.Request.Context(), service.RoomCreateInput{
		Name:   inp.Name,
		UserId: userId,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, roomCreateResponse{
		ID:    roomId,
		Token: token,
	})
}
