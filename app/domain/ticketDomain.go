package domain

type AddTicketPayload struct {
	Issuer_code      string `json:"issuer_code" validation:"required"`
	Ticket_ref       string `json:"ticket_ref" validation:"required"`
	Passenger_id     string `json:"passenger_id" validation:"required"`
	Passenger_name   string `json:"passenger_name" validation:"required"`
	Passenger_msisdn string `json:"passenger_msisdn" validation:"required"`
	Origin_code      string `json:"origin_code" validation:"required"`
	Destination_code string `json:"destination_code" validation:"required"`
	Booking_at       string `json:"booking_at" validation:"required"`
	Issuer_fare      int    `json:"issuer_fare" validation:"required"`
	Ticket_group     string `json:"ticket_group,omitempty"`
}

type AddTicketResponse struct {
	Ref_no      string `json:"ref_no"`
	Ticket_code string `json:"ticket_code"`
}

type DeletePayload struct {
	Ticket_code string `json:"ticket_code" param:"ticket_code" validation:"required"`
}

type SPDelete struct {
	Status_code string `json:"status_code"`
	Message     string `json:"message"`
	Ref_no      string `json:"ref_no,omitempty"`
}

type DataTypeAddTicket struct {
	Status_code string `json:"status_code"`
	Message     string `json:"message"`
	Ref_no      string `json:"ref_no"`
	Ticket_code string `json:"ticket_code,omitempty"`
}

type GetTicketByTicketCode struct {
	Fare             int    `json:"fare"`
	Issuer_fare      int    `json:"issuer_dare"`
	Operational_date string `json:"operational_date"`
}