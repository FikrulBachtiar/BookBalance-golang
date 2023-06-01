package domain

type PayloadFare struct {
	Origin      string `json:"origin" param:"origin" validate:"required"`
	Destination string `json:"destination" param:"destination" validate:"required"`
}