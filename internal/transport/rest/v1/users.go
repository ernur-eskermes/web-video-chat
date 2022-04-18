package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.GET("/me", h.userGetMe)
			authenticated.POST("/subscriptions/create", h.createSubscription)
		}
	}
}

type userAccountResponse struct {
	ID       primitive.ObjectID   `json:"id"`
	Username string               `json:"username"`
	ACCSubs  []primitive.ObjectID `json:"acc_subs" bson:"acc_subs"`
	PNDSubs  []primitive.ObjectID `json:"pnd_subs" bson:"pnd_subs"`
}

// @Summary Get User Info
// @Tags users
// @Security UsersAuth
// @Description user get me
// @Accept  json
// @Produce  json
// @Success 200 {object} userAccountResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/me [get]
func (h *Handler) userGetMe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user, err := h.services.Users.GetById(c.Request.Context(), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, userAccountResponse{
		ID:       userId,
		Username: user.Username,
		ACCSubs:  user.ACCSubs,
		PNDSubs:  user.PNDSubs,
	})
}

type createSubscriptionInput struct {
	UserId primitive.ObjectID `json:"user_id"`
}

// @Summary Create User Subscription
// @Tags users
// @Security UsersAuth
// @Description create user subscription
// @Accept  json
// @Produce  json
// @Param input body createSubscriptionInput true "create subscription"
// @Success 201 {string} string ok
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/subscriptions/create [post]
func (h *Handler) createSubscription(c *gin.Context) {
	var inp createSubscriptionInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Users.CreateSubscription(c, userId, inp.UserId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusCreated)
}
