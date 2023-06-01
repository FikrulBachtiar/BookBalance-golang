package controllers

import (
	"bookbalance/app/configs"
	"bookbalance/app/domain"
	"bookbalance/app/services"
	"bookbalance/app/utils"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TicketController interface {
	AddTicket(ctx echo.Context) error
}

type ticketController struct {
	ticketService services.TicketService
	db *sql.DB
}

func NewTicketController(e *echo.Echo, ticketService services.TicketService, db *sql.DB) TicketController {
	return &ticketController{
		ticketService: ticketService,
		db: db,
	}
}

func (tc *ticketController) AddTicket(ctx echo.Context) error {
	
	var dataResponse domain.AddTicketResponse 
	payload := new(domain.AddTicketPayload);

	if err := ctx.Bind(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusBadRequest,
			Code: 10,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if err := ctx.Validate(payload); err != nil {
		response := &configs.Response{
			Status: http.StatusInternalServerError,
			Code: 11,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	// Check booking_at format RFC 3339
	// if _, err := time.Parse(time.RFC3339, payload.Booking_at); err != nil {
	// 	response := &configs.Response{
	// 		Status: http.StatusBadRequest,
	// 		Code: 67,
	// 		Message: "Invalid time format. Use the RFC 3339 format",
	// 	}
	// 	return response.ResponseMiddleware(ctx);
	// }

	// fmt.Println("Damn => ", payload.Issuer_code);

	payload.Passenger_msisdn = utils.TrimPhoneNumber(payload.Passenger_msisdn);

	status, code, err := tc.ticketService.AddTicket(ctx.Request().Context(), payload, &dataResponse);
	if err != nil {
		response := &configs.Response{
			Status: status,
			Code: code,
			Error: err.Error(),
		}
		return response.ResponseMiddleware(ctx);
	}

	if code != 0 {
		response := &configs.Response{
			Status: status,
			Code: code,
		}
		return response.ResponseMiddleware(ctx);
	}

	response := &configs.Response{
		Status: status,
		Code: code,
		Data: &dataResponse,
	}
	return response.ResponseMiddleware(ctx);

}