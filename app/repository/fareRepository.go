package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type FareRepository interface {
	GetStationIssuer(ctx context.Context, issuer_code string, origin_code string, destination_code string) ([]string, error)
	GetFareData(ctx context.Context, origin_code string, destination_code string) (int, error)
}

type fareRepository struct {
	db *sql.DB
}

func NewFareRepository(db *sql.DB) FareRepository {
	return &fareRepository{
		db: db,
	}
}

func (fr *fareRepository) GetStationIssuer(ctx context.Context, issuer_code string, origin_code string, destination_code string) ([]string, error) {
	var list []string;
	sqlQuery := fmt.Sprintf("SELECT station_code FROM mtr.v_station_issuer WHERE station_code IN('%s', '%s') AND third_party_code = '%s' AND active_status = '1'", origin_code, destination_code, issuer_code);
	
	rows, err := fr.db.Query(sqlQuery);
	if err != nil {
		return list, err;
	}
	defer rows.Close();

	for rows.Next() {
		var station_code string;
		err = rows.Scan(&station_code);
		if err != nil {
			return list, err;
		}
		
		list = append(list, station_code)
	}

	return list, nil;
}

func (fr *fareRepository) GetFareData(ctx context.Context, origin_code string, destination_code string) (int, error) {
	sqlQuery := fmt.Sprintf("SELECT fare FROM mtr.t_mtr_fare WHERE origin_code = '%s' AND destination_code = '%s' AND active_status = '1'", origin_code, destination_code);
	var fare int

	err := fr.db.QueryRowContext(ctx, sqlQuery).Scan(&fare);
	if err != nil {
		return 0, err;
	}

	return fare, nil;
}