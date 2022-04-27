package v1

import (
	"errors"

	"github.com/ernur-eskermes/web-video-chat/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
		}

		c.AbortWithStatusJSON(statusCode, gin.H{"errors": out})

		return
	}

	logger.Error(err.Error())
	c.AbortWithStatusJSON(statusCode, response{err.Error()})
}
