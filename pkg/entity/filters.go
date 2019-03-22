package entity

import (
	"fmt"
	"time"
)

// Filters contains a filter's information.
type Filters struct {
	DriverID             string    `schema:"driverId,ommitempty"`
	Seats                *int      `schema:"seats,ommitempty"`
	LeaveAt              time.Time `schema:"leaveAt,ommitempty"`
	ArriveBy             time.Time `schema:"arriveBy,ommitempty"`
	DetailsAnimals       *int      `schema:"detailsAnimals,ommitempty"`
	DetailsLuggages      *int      `schema:"detailsLuggages,ommitempty"`
	RadiusThresh         *int      `schema:"radiusThresh,ommitempty"`
	DestinationLatitude  *float64  `schema:"destinationLatitude,ommitempty"`
	DestinationLongitude *float64  `schema:"destinationLongitude,ommitempty"`
}

const (
	// MinimumLuggagesValue represents the minimum value for luggages
	MinimumLuggagesValue = 0

	// MaximumLuggagesValue represents the maximum value for luggages
	MaximumLuggagesValue = 2

	// MinimumAnimalsValue represents the minimum value for animals
	MinimumAnimalsValue = 0

	// MaximumAnimalsValue represents the maximum value for animals
	MaximumAnimalsValue = 1

	// MinimumRadiusThresh represents the minimum value for radius threshold
	MinimumRadiusThresh = 0
)

// Validate validates that the filters's required fields are filled out correctly.
func (f *Filters) Validate() error {
	if f.Seats != nil && *f.Seats < 0 {
		return ValidationError{"seats filter must be greater than 0"}
	}

	if !f.LeaveAt.IsZero() && !f.ArriveBy.IsZero() {
		return ValidationError{"can't have leaveAt and arriveBy filter at the same time"}
	}

	if f.DetailsAnimals != nil && (*f.DetailsAnimals < MinimumAnimalsValue || *f.DetailsAnimals > MaximumAnimalsValue) {
		return ValidationError{fmt.Sprintf("detailsAnimals filter must be between ")}
	}

	if f.DetailsLuggages != nil && (*f.DetailsLuggages < MinimumLuggagesValue || *f.DetailsLuggages > MaximumLuggagesValue) {
		return ValidationError{fmt.Sprintf("detailsLuggages filter must be between %d and %d", MinimumLuggagesValue, MaximumLuggagesValue)}
	}

	if f.RadiusThresh != nil && *f.RadiusThresh <= MinimumRadiusThresh {
		return ValidationError{fmt.Sprintf("radiusThresh must be greater than %d", MinimumRadiusThresh)}
	}

	if f.DestinationLongitude != nil && (*f.DestinationLongitude < MinimumLongitude || *f.DestinationLongitude > MaximumLongitude) {
		return ValidationError{"invalid destination longitude value"}
	}

	if f.DestinationLatitude != nil && (*f.DestinationLatitude < MinimumLatitude || *f.DestinationLatitude > MaximumLatitude) {
		return ValidationError{"invalid destination latitude value"}
	}

	return nil
}
