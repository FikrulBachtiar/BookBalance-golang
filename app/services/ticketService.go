package services

import (
	"bookbalance/app/domain"
	"bookbalance/app/repository"
	"bookbalance/app/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TicketService interface {
	AddTicket(ctx context.Context, data *domain.AddTicketPayload, res *domain.AddTicketResponse) (int, int, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepo
}

func NewTicketService(ticketRepo repository.TicketRepo) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
	}
}

func (ts *ticketService) AddTicket(ctx context.Context, data *domain.AddTicketPayload, res *domain.AddTicketResponse) (int, int, error) {

	fmt.Println("DATA => ", data.Issuer_code)
	err := ts.ticketRepo.GetIssuerByCode(ctx, data.Issuer_code);
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusBadRequest, 23, nil;
		} else {
			return http.StatusInternalServerError, 500, err;
		}
	}

	stations, err := ts.ticketRepo.GetStationPermitted(ctx, data.Issuer_code, data.Origin_code, data.Destination_code); 
	if err != nil {
		return http.StatusInternalServerError, 500, err;
	}

	originValid := utils.StationContains(stations, data.Origin_code);
	destinationValid := utils.StationContains(stations, data.Destination_code);

	if !originValid {
		return http.StatusBadRequest, 23, nil;
	}

	if !destinationValid {
		return http.StatusBadRequest, 24, nil;
	}

	result, err := ts.ticketRepo.SPAddTicket(ctx, data);
	if err != nil {
		return http.StatusInternalServerError, 500, err;
	}

	num, err := strconv.Atoi(result.Status_code);
	if err != nil {
		return http.StatusInternalServerError, 500, err;
	}
	
	tickets, err := ts.ticketRepo.GetTicketByTicketCode(result.Ticket_code);
	if err != nil {
		return http.StatusInternalServerError, 500, err;
	}

	date, err := time.Parse("2006-01-02T00:00:00Z", tickets.Operational_date);
	if err != nil {
		return http.StatusInternalServerError, 500, err;
	}
	dateFormat := date.Format("2006-01-02");

	res.Ref_no = result.Ref_no;
	res.Ticket_code = result.Ticket_code;

	go ts.ticketRepo.IncrementSummary(ctx, data.Issuer_code, data.Origin_code, data.Destination_code, "vol", 1, dateFormat);
	go ts.ticketRepo.IncrementSummary(ctx, data.Issuer_code, data.Origin_code, data.Destination_code, "rev", tickets.Fare, dateFormat);
	go ts.ticketRepo.IncrementSummary(ctx, data.Issuer_code, data.Origin_code, data.Destination_code, "iss", tickets.Issuer_fare, dateFormat);

	return http.StatusOK, num, nil;
}