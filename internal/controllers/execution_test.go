package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"quantito-engine/internal/interfaces/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		router := echo.New()

		id := 1

		service := mocks.NewMockIExecutionService(t)

		service.EXPECT().Execute(id).
			Once().
			Return(nil)

		controller := NewExecutionController(service)

		request := httptest.NewRequest(http.MethodPost, "/execute", nil)
		responseRecorder := httptest.NewRecorder()
		c := router.NewContext(request, responseRecorder)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		assert.NoError(t, controller.Execute(c))
		assert.Equal(t, http.StatusOK, responseRecorder.Code)
	})

	t.Run("Returns HTTP 400 if id param is not a valid int", func(t *testing.T) {
		router := echo.New()

		service := mocks.NewMockIExecutionService(t)
		controller := NewExecutionController(service)

		request := httptest.NewRequest(http.MethodPost, "/execute", nil)
		responseRecorder := httptest.NewRecorder()
		c := router.NewContext(request, responseRecorder)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid_id")

		assert.NoError(t, controller.Execute(c))
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	})

	t.Run("Returns 422 if an error occurred in lower layers", func(t *testing.T) {
		router := echo.New()

		id := 1

		service := mocks.NewMockIExecutionService(t)

		service.EXPECT().Execute(id).
			Once().
			Return(errors.New("ERROR"))

		controller := NewExecutionController(service)

		request := httptest.NewRequest(http.MethodPost, "/execute", nil)
		responseRecorder := httptest.NewRecorder()
		c := router.NewContext(request, responseRecorder)

		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		assert.NoError(t, controller.Execute(c))
		assert.Equal(t, http.StatusUnprocessableEntity, responseRecorder.Code)

		responseBody := map[string]any{}
		err := json.NewDecoder(responseRecorder.Body).Decode(&responseBody)

		assert.NoError(t, err)
		assert.Equal(t, "ERROR", responseBody["error"])
	})
}
