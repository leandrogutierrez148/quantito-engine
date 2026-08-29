package interfaces

import (
	"github.com/labstack/echo/v4"
)

type IExecutionController interface {
	Execute(c echo.Context) error
}

type IExecutionService interface {
	Execute(id int) error
}
