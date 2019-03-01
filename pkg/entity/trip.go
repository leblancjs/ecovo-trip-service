package entity

import (
	"fmt"
	"time"
)

// Trip contains a trips's information.
type Trip struct {
	ID          ID        `json:"id"`
	DriverID    ID        `json:"driverId"`
	VehicleID   ID        `json:"vehicleId"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	LeaveAt     time.Time `json:"leaveAt"`
	ArriveBy    time.Time `json:"arriveBy"`
	Seats       int       `json:"seats"`
	Stops       []string  `json:"stops"`
	Details     *Details  `json:"details"`
}

const (
	// MinimumSeats represents the minimum seats possible in a car.
	MinimumSeats = 1

	// MaximumSeats represents the maximum seats possible in a car.
	MaximumSeats = 10
)

// Validate validates that the trips's required fields are filled out correctly.
func (t *Trip) Validate() error {
	if t.Source == "" {
		return ValidationError{"source is missing"}
	}

	if t.Destination == "" {
		return ValidationError{"destination is missing"}
	}

	if t.Seats < MinimumSeats || t.Seats > MaximumSeats {
		return ValidationError{fmt.Sprintf("number of seats must be between %d and %d", MinimumSeats, MaximumSeats)}
	}

	if t.Details != nil {
		err := t.Details.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
