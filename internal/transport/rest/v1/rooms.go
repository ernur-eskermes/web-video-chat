package v1

import (
	"errors"
	"net/http"

	"github.com/ernur-eskermes/web-video-chat/internal/core"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) initRoomsRoutes(api *gin.RouterGroup) {
	rooms := api.Group("/rooms")
	{
		authenticated := rooms.Group("/", h.userIdentity)
		{
			authenticated.GET("", h.roomList)
			authenticated.POST("/create", h.roomCreate)
			authenticated.POST("/:roomID/join", h.roomJoin)
		}
	}
}

type roomCreateResponse struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}

// @Summary Room Create
// @Tags rooms
// @Security UsersAuth
// @Description create room
// @ModuleID roomCreate
// @Accept  json
// @Produce  json
// @Param input body core.RoomCreateInput true "sign up info"
// @Success 201 {object} roomCreateResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /rooms/create [post]
func (h *Handler) roomCreate(c *gin.Context) {
	var inp core.RoomCreateInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	roomId, token, err := h.services.Rooms.Create(c.Request.Context(), inp, userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, roomCreateResponse{
		ID:    roomId,
		Token: token,
	})
}

// @Summary Room List
// @Tags rooms
// @Security UsersAuth
// @Description get rooms
// @ModuleID roomList
// @Accept  json
// @Produce  json
// @Success 200 {object} core.Room
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /rooms [get]
func (h *Handler) roomList(c *gin.Context) {
	if _, err := getUserId(c); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	rooms, err := h.services.Rooms.GetList(c.Request.Context(), true)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, rooms)
}

type roomTokenResponse struct {
	Token string `json:"token"`
}

// @Summary Join Room
// @Tags rooms
// @Security UsersAuth
// @Description Join Room
// @ModuleID roomJoin
// @Accept  json
// @Produce  json
// @Param roomID path string true "room id"
// @Success 200 {object} roomTokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /rooms/{roomID}/join [post]
func (h *Handler) roomJoin(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	roomID, err := getIdByContext(c, "roomID")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	token, err := h.services.Rooms.GetByID(c.Request.Context(), roomID, userID)
	if err != nil {
		if errors.Is(err, core.ErrRoomNotFound) {
			newResponse(c, http.StatusNotFound, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, roomTokenResponse{
		Token: token,
	})
}
