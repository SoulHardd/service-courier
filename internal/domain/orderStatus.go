package domain

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)
