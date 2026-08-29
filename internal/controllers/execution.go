package controllers

import (
	"net/http"
	"strconv"

	"quantito-engine/internal/interfaces"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ExecutionController handles HTTP execution requests for experiments.
type ExecutionController struct {
	service interfaces.IExecutionService
}

// NewExecutionController instantiates a new ExecutionController with its dependencies.
func NewExecutionController(
	service interfaces.IExecutionService,
) *ExecutionController {
	return &ExecutionController{
		service: service,
	}
}

type executionErrorResponse struct {
	Error string `json:"error"`
}

// Execute extracts the experiment ID parameter and dispatches the execution to the service.
func (controller *ExecutionController) Execute(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, executionErrorResponse{"Invalid ID parameter"})
	}

	err = controller.service.Execute(id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Error executing experiment")
		return c.JSON(http.StatusUnprocessableEntity, executionErrorResponse{err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
