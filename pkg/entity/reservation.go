package entity

import (
	"fmt"
)

// Reservation contains a reservation's information.
type Reservation struct {
	TripID        ID  `json:"tripId"`
	UserID        ID  `json:"userId"`
	SourceID      ID  `json:"sourceId"`
	DestinationID ID  `json:"destinationId"`
	Seats         int `json:"seats"`
}

// Validate validates that the reservation's required fields are filled out correctly.
func (r *Reservation) Validate() error {
	if r.TripID.IsZero() {
		return ValidationError{"Trip's ID is missing"}
	}

	if r.UserID.IsZero() {
		return ValidationError{"User's ID is missing"}
	}

	if r.SourceID.IsZero() {
		return ValidationError{"Source's ID is missing"}
	}

	if r.DestinationID.IsZero() {
		return ValidationError{"Destination's ID is missing"}
	}

	if r.Seats < MinimumSeats || r.Seats > MaximumSeats {
		return ValidationError{fmt.Sprintf("number of seats must be between %d and %d", MinimumSeats, MaximumSeats)}
	}

	return nil
}
