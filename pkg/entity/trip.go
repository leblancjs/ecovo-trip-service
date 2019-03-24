package entity

import (
	"fmt"
	"time"
)

// Trip contains a trips's information.
type Trip struct {
	ID        ID        `json:"id"`
	DriverID  ID        `json:"driverId"`
	VehicleID ID        `json:"vehicleId"`
	Full      bool      `json:"full"`
	LeaveAt   time.Time `json:"leaveAt"`
	ArriveBy  time.Time `json:"arriveBy"`
	Seats     int       `json:"seats"`
	Stops     []*Stop   `json:"stops"`
	Details   *Details  `json:"details"`
}

const (
	// MinimumSeats represents the minimum seats possible in a car.
	MinimumSeats = 1

	// MaximumSeats represents the maximum seats possible in a car.
	MaximumSeats = 10
)

// Validate validates that the trips's required fields are filled out correctly.
func (t *Trip) Validate() error {
	if t.LeaveAt.IsZero() && t.ArriveBy.IsZero() {
		return ValidationError{"leaveAt or arriveBy is missing"}
	}

	if !t.LeaveAt.IsZero() && t.LeaveAt.Before(time.Now()) {
		return ValidationError{"leaveAt must be in the future"}
	}

	if !t.ArriveBy.IsZero() && t.ArriveBy.Before(time.Now()) {
		return ValidationError{"arriveBy must be in the future"}
	}

	if t.DriverID.IsZero() {
		return ValidationError{"Driver's ID is missing"}
	}

	if t.VehicleID.IsZero() {
		return ValidationError{"Vehicle's ID is missing"}
	}

	if t.Seats < MinimumSeats || t.Seats > MaximumSeats {
		return ValidationError{fmt.Sprintf("number of seats must be between %d and %d", MinimumSeats, MaximumSeats)}
	}

	if t.Details != nil {
		err := t.Details.Validate()
		if err != nil {
			return err
		}
	} else {
		return ValidationError{"missing details"}
	}

	if t.Stops != nil {
		for _, s := range t.Stops {
			err := s.Validate()
			if err != nil {
				return err
			}
		}
	} else {
		return ValidationError{"missing stops"}
	}

	return nil
}
