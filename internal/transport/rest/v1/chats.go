package v1

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ernur-eskermes/web-video-chat/internal/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func (h *Handler) initChatsRoutes(api *gin.RouterGroup) {
	chats := api.Group("/chats")
	{
		chats.GET("/:id", h.indexChat)
		chats.GET("/ws/:id", h.chatWebsocket)
	}
}

func (h *Handler) indexChat(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) chatWebsocket(c *gin.Context) {
	chatId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	lock := new(sync.Mutex)

	h.websocket.HandleMessage(func(s *melody.Session, msg []byte) {
		_ = h.websocket.Broadcast(msg)

		_ = h.services.Chats.CreateMessage(c, service.CreateMessageInput{
			Message: string(msg),
			UserId:  userId,
			ChatId:  chatId,
		})
	})
	h.websocket.HandleConnect(func(s *melody.Session) {
		lock.Lock()
		messages, _ := h.services.Chats.GetMessages(c, chatId)
		b, _ := json.Marshal(messages)
		s.Write(b)
		lock.Unlock()
	})

	if err = h.websocket.HandleRequest(c.Writer, c.Request); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}
}
