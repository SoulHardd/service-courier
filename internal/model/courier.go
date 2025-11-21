package model

import "time"

type Courier struct {
	ID        int64
	Name      string
	Phone     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
