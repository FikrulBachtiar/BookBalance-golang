package repository

import (
	"bookbalance/app/domain"
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type TicketRepo interface {
	SPDeleteTicket(ctx context.Context, issuer_code string, ticket_code string) error
	GetIssuerByCode(ctx context.Context, issuer_code string) error
	GetStationPermitted(ctx context.Context, issuer_code string, origin_code string, destination_code string) ([]string, error)
	SPAddTicket(ctx context.Context, data *domain.AddTicketPayload) (*domain.DataTypeAddTicket, error)
	// IncrementSummary(ctx context.Context, issuer_code string, origin string, destination string, types string, increment int, summary_date string) error
	IncrementSummary(ctx context.Context, issuer_code string, origin string, destination string, types string, increment int, summary_date string)
	GetTicketByTicketCode(ticket_code string) (domain.GetTicketByTicketCode, error)
}

type ticketRepo struct {
	db *sql.DB
	redis *redis.Client
}

func NewTicketRepo(db *sql.DB, redis *redis.Client) TicketRepo {
	return &ticketRepo{
		db: db,
		redis: redis,
	}
}

func (tr *ticketRepo) SPDeleteTicket(ctx context.Context, issuer_code string, ticket_code string) error {

	sqlQuery := fmt.Sprintf("SELECT * FROM trx.sp_trx_ticket_delete('%s', '%s')", issuer_code, ticket_code);
	var result domain.SPDelete

	err := tr.db.QueryRow(sqlQuery).Scan(&result);
	if err != nil {
		return err;
	}

	fmt.Println("Data => ", result)

	return nil;
}

func (tr *ticketRepo) GetIssuerByCode(ctx context.Context, issuer_code string) error {

	sqlQuery := fmt.Sprintf("SELECT third_party_code FROM mtr.t_mtr_third_party WHERE third_party_code = '%s' AND active_status = 1", issuer_code);
	var code string

	err := tr.db.QueryRow(sqlQuery).Scan(&code);
	if err != nil {
		return err;
	}

	return nil;
}

func (tr *ticketRepo) GetStationPermitted(ctx context.Context, issuer_code string, origin_code string, destination_code string) ([]string, error) {

	sqlQuery := fmt.Sprintf("SELECT station_code FROM mtr.v_station_issuer WHERE station_code IN('%s', '%s') AND third_party_code = '%s' AND active_status = '1'", origin_code, destination_code, issuer_code);
	var stations []string;

	rows, err := tr.db.Query(sqlQuery);
	if err != nil {
		return make([]string, 0), err;
	}
	defer rows.Close();
	
	for rows.Next() {
		var station string

		err = rows.Scan(&station);
		if err != nil {
			return make([]string, 0), err;
		}

		stations = append(stations, station);
	}

	return stations, err;
}

func (tr *ticketRepo) SPAddTicket(ctx context.Context, data *domain.AddTicketPayload) (*domain.DataTypeAddTicket, error) {
	
	sqlQuery := fmt.Sprintf("SELECT * FROM trx.sp_trx_ticket_add_v3('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'::TIMESTAMP, %d, '%s')",
		data.Issuer_code, data.Ticket_ref, data.Passenger_id, data.Passenger_name, data.Passenger_msisdn, data.Origin_code, data.Destination_code, 
		data.Booking_at, data.Issuer_fare, data.Ticket_group);
	var dataType domain.DataTypeAddTicket;

	err := tr.db.QueryRow(sqlQuery).Scan(&dataType.Status_code, &dataType.Message, &dataType.Ref_no, &dataType.Ticket_code);
	if err != nil {
		return &dataType, err;
	}

	return &dataType, nil;
}

func (tr *ticketRepo) IncrementSummary(ctx context.Context, issuer_code string, origin string, destination string, types string, increment int, summary_date string) {
	key := fmt.Sprintf("summary:%s", summary_date);
	field := fmt.Sprintf("%s:%s-%s:%s", issuer_code, origin, destination, types);

	ctxRedis := context.Background();
	if tr.redis != nil {
		ctxRedis = tr.redis.Context()
	}

	tr.redis.HIncrBy(ctxRedis, key, field, int64(increment));
	// if err != nil {
	// 	fmt.Println("HELLO")
	// 	log.Fatal(err);
	// }
}

func (tr *ticketRepo) GetTicketByTicketCode(ticket_code string) (domain.GetTicketByTicketCode, error) {
	sqlQuery := fmt.Sprintf("SELECT fare, issuer_fare, operational_date FROM trx.t_trx_ticket WHERE ticket_code = '%s'", ticket_code);
	var tickets domain.GetTicketByTicketCode;

	err := tr.db.QueryRow(sqlQuery).Scan(&tickets.Fare, &tickets.Issuer_fare, &tickets.Operational_date);
	if err != nil {
		return tickets, err;
	}

	return tickets, nil;
}