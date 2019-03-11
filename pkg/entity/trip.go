package entity

import (
	"fmt"
	"time"
)

// Trip contains a trips's information.
type Trip struct {
	ID          ID        `json:"id"`
	Driver      *Driver   `json:"driver"`
	Vehicle     *Vehicle  `json:"vehicle"`
	Source      *Map      `json:"source"`
	Destination *Map      `json:"destination"`
	LeaveAt     time.Time `json:"leaveAt"`
	ArriveBy    time.Time `json:"arriveBy"`
	Seats       int       `json:"seats"`
	Stops       []*Map    `json:"stops"`
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
	if t.LeaveAt.IsZero() && t.ArriveBy.IsZero() {
		return ValidationError{"leaveAt or arriveBy is missing"}
	}

	if !t.LeaveAt.IsZero() && !t.ArriveBy.IsZero() {
		return ValidationError{"can't have leaveAt and arriveBy"}
	}

	if t.Driver.ID.IsZero() {
		return ValidationError{"Driver's ID is missing"}
	}

	if t.Vehicle.ID.IsZero() {
		return ValidationError{"Vehicle's ID is missing"}
	}

	if t.Seats < MinimumSeats || t.Seats > MaximumSeats {
		return ValidationError{fmt.Sprintf("number of seats must be between %d and %d", MinimumSeats, MaximumSeats)}
	}

	if t.Source != nil {
		err := t.Source.Validate()
		if err != nil {
			return err
		}
	}

	if t.Destination != nil {
		err := t.Destination.Validate()
		if err != nil {
			return err
		}
	}

	if t.Driver != nil {
		err := t.Driver.Validate()
		if err != nil {
			return err
		}
	}

	if t.Details != nil {
		err := t.Details.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
