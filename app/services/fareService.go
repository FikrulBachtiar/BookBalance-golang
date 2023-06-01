package services

import (
	"bookbalance/app/repository"
	"bookbalance/app/utils"
	"context"
)

type FareService interface {
	IsStationPermitted(ctx context.Context, issuer_code string, origin_code string, destination_code string) (int, error)
	GetFare(ctx context.Context, origin_code string, destination_code string) (int, int, error)
}

type fareService struct {
	fareRepo repository.FareRepository
}

func NewFareService(fareRepo repository.FareRepository) FareService {
	return &fareService{
		fareRepo: fareRepo,
	}
}

func (fs *fareService) IsStationPermitted(ctx context.Context, issuer_code string, origin_code string, destination_code string) (int, error) {
	list_station, err := fs.fareRepo.GetStationIssuer(ctx, issuer_code, origin_code, destination_code);
	if err != nil {
		return 500, err;
	}

	if (len(list_station) <= 0) {
		return 88, nil;
	}

	originValid := utils.StationContains(list_station, origin_code);
	destinationValid := utils.StationContains(list_station, destination_code);

	if !originValid {
		return 23, nil;
	}

	if !destinationValid {
		return 24, nil;
	}

	return 0, nil;
}

func (fs *fareService) GetFare(ctx context.Context, origin_code string, destination_code string) (int, int, error) {
	fare, err := fs.fareRepo.GetFareData(ctx, origin_code, destination_code);
	if err != nil {
		// Fare not found
		return 60, 0, err;
	}

	return 0, fare, nil;
}