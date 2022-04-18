package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/olahol/melody.v1"

	"github.com/ernur-eskermes/web-video-chat/internal/config"
	"github.com/ernur-eskermes/web-video-chat/internal/service"
	handler "github.com/ernur-eskermes/web-video-chat/internal/transport/rest"
	"github.com/ernur-eskermes/web-video-chat/pkg/auth"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	h := handler.NewHandler(&service.Services{}, &auth.Manager{}, &melody.Melody{})

	require.IsType(t, &handler.Handler{}, h)
}

func TestNewHandler_Init(t *testing.T) {
	h := handler.NewHandler(&service.Services{}, &auth.Manager{}, &melody.Melody{})

	router := h.Init(&config.Config{
		Limiter: config.LimiterConfig{
			RPS:   2,
			Burst: 4,
			TTL:   10 * time.Minute,
		},
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	require.Equal(t, http.StatusOK, res.StatusCode)
}
