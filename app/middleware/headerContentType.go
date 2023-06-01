package middleware

import (
	"github.com/labstack/echo/v4"
)

func ContentTypeResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON);
		return next(c)
	}
}