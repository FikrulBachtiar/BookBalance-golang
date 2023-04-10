package routes

import (
	"bookbalance/model"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func InitRoutes() *echo.Echo {
	app := echo.New()
	app.Validator = &model.PayloadValidator{Validator: validator.New()}

	return app
}