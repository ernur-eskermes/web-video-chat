package v1

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ernur-eskermes/web-video-chat/internal/domain"
	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.GET("/me", h.userGetMe)
		}

		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/auth/refresh", h.userRefresh)
	}
}

type userSignUpInput struct {
	Username        string `json:"username" binding:"required,max=64"`
	Password        string `json:"password" binding:"required,min=8,max=64"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,max=64"`
}

type signInInput struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-up [post]
func (h *Handler) userSignUp(c *gin.Context) {
	var inp userSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
		Username:        inp.Username,
		Password:        inp.Password,
		ConfirmPassword: inp.ConfirmPassword,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusCreated)
}

// @Summary User SignIn
// @Tags users-auth
// @Description user sign in
// @ModuleID userSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-in [post]
func (h *Handler) userSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
		Username: inp.Username,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/auth/refresh [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.RefreshTokens(c.Request.Context(), inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type userAccountResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

// @Summary Get User Info
// @Tags users-auth
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
	})
}
