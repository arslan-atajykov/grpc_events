package order

import "time"

// Order is our domain model (internal struct).
type Order struct {
	ID        int64
	Customer  string
	Status    string
	CreatedAt time.Time
}
