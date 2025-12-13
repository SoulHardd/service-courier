package domain

type CourierStatus string

const (
	CourierStatusAvailable CourierStatus = "available"
	CourierStatusBusy      CourierStatus = "busy"
)
