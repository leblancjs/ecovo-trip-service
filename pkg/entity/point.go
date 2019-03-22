package entity

import (
	"fmt"
)

// Point contains a geolocation's information.
type Point struct {
	Longitude float64 `json:"longitude" schema:"longitude"`
	Latitude  float64 `json:"latitude" schema:"latitude"`
	Name      string  `json:"name" schema:"name,ommitempty"`
}

const (
	// MinimumLongitude represents the minimum longitude value.
	MinimumLongitude = -180

	// MaximumLongitude represents the maximum longitude value.
	MaximumLongitude = 180

	// MinimumLatitude represents the minimum latitude value.
	MinimumLatitude = -90

	// MaximumLatitude represents the maximum latitude value.
	MaximumLatitude = 90
)

// String returns string value of Point.
func (p *Point) String() string {
	return fmt.Sprintf("%f", p.Latitude) + ", " + fmt.Sprintf("%f", p.Longitude)
}

// Validate validates that the map's required fields are filled out correctly.
func (p *Point) Validate() error {
	if p.Longitude < MinimumLongitude || p.Longitude > MaximumLongitude {
		return ValidationError{"invalid longitude value"}
	}

	if p.Latitude < MinimumLatitude || p.Latitude > MaximumLatitude {
		return ValidationError{"invalid latitude value"}
	}

	if p.Name == "" {
		return ValidationError{"name is empty"}
	}

	return nil
}
