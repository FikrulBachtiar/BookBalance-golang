package controllers

import (
	"bookbalance/app/configs"
	"bookbalance/app/domain"
	"bookbalance/app/services"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type fareController struct {
	fareService services.FareService
	db *sql.DB
}

func NewFareController(e *echo.Echo, fareService services.FareService, db *sql.DB) *fareController {
	return &fareController{
		fareService: fareService,
		db: db,
	}
}

func (fc *fareController) Fare(ctx echo.Context) error {

	payload := new(domain.PayloadFare)
	issuer_code := ctx.QueryParam("issuer_code")

	// Binding payload (Request) into struct
	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 10,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	// Payload validation to match the struct
	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: 11,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	code, err := fc.fareService.IsStationPermitted(ctx.Request().Context(), issuer_code, payload.Origin, payload.Destination);
	if err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: code,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Code: code,
			Message: "data not found",
		}
		return response.ResponseMiddleware(ctx);
	}

	code, fare, err := fc.fareService.GetFare(ctx.Request().Context(), payload.Origin, payload.Destination);
	if err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: code,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Code: code,
			Message: "data not found",
		}
		return response.ResponseMiddleware(ctx);
	}
	
	response := &configs.Response{
		Status: http.StatusOK,
		Code: 0,
		Message: "data found",
		Data: map[string]interface{}{"fare": fare},
	}
	return response.ResponseMiddleware(ctx);
}