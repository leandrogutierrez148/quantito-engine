package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"quantito-engine/internal/controllers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	controller := controllers.NewPingController()

	router := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	responseRecorder := httptest.NewRecorder()
	c := router.NewContext(request, responseRecorder)

	assert.NoError(t, controller.Ping(c))
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, "Pong", responseRecorder.Body.String())
}
