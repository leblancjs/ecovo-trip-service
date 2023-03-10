package route

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"azure.com/ecovo/trip-service/pkg/entity"
	"googlemaps.github.io/maps"
)

// GoogleMapsRepository structure
type GoogleMapsRepository struct {
	client *maps.Client
}

const (
	// PricePerKilometer represents the price (in dollars) established by the Canadian
	// government as the maximum allocation for vehicle usage.
	// https://www.canada.ca/fr/agence-revenu/services/impot/entreprises/sujets/retenues-paie/avantages-allocations/automobile/allocations-frais-automobile-vehicule-a-moteur/taux-allocations-frais-automobile.html
	PricePerKilometer = 0.58
)

// NewGoogleMapsRepository creates the repository
func NewGoogleMapsRepository(client *maps.Client) (Repository, error) {
	if client == nil {
		return nil, fmt.Errorf("route.GoogleRepository: client is nil")
	}

	return &GoogleMapsRepository{
		client: client,
	}, nil
}

// GenerateRoute allows a user to generate route information for a given trip
func (gr *GoogleMapsRepository) GenerateRoute(t *entity.Trip) error {
	var wp = make([]string, len(t.Stops))
	for i, s := range t.Stops {
		wp[i] = s.Point.String()
	}

	var dr *maps.DirectionsRequest

	if t.LeaveAt.IsZero() && !t.ArriveBy.IsZero() {
		dr = &maps.DirectionsRequest{
			Origin:      t.Stops[0].Point.String(),
			Destination: t.Stops[len(t.Stops)-1].Point.String(),
			Waypoints:   wp,
			ArrivalTime: strconv.FormatInt(t.ArriveBy.Unix(), 10),
		}
	} else if !t.LeaveAt.IsZero() && t.ArriveBy.IsZero() {
		dr = &maps.DirectionsRequest{
			Origin:        t.Stops[0].Point.String(),
			Destination:   t.Stops[len(t.Stops)-1].Point.String(),
			Waypoints:     wp,
			DepartureTime: strconv.FormatInt(t.LeaveAt.Unix(), 10),
		}
	} else if !t.LeaveAt.IsZero() && !t.ArriveBy.IsZero() {
		dr = &maps.DirectionsRequest{
			Origin:        t.Stops[0].Point.String(),
			Destination:   t.Stops[len(t.Stops)-1].Point.String(),
			Waypoints:     wp,
			DepartureTime: strconv.FormatInt(t.LeaveAt.Unix(), 10),
		}
	} else {
		return fmt.Errorf("trip.GoogleMapsRepository: arriveBy OR leaveAt must be specified")
	}

	r, _, err := gr.client.Directions(context.Background(), dr)
	if err != nil {
		return fmt.Errorf("trip.GoogleMapsRepository: error getting directions, %s", err)
	}

	// Here we modify our trip data with google maps route generated
	if len(r) > 0 {
		route := r[0]

		if t.LeaveAt.IsZero() {
			leaveAt := t.ArriveBy
			for i := range route.Legs {
				t.TotalDistance += route.Legs[i].Distance.Meters
				leaveAt = leaveAt.Add(-(route.Legs[i].Duration) * time.Nanosecond)
			}
			t.LeaveAt = leaveAt
		}

		if t.ArriveBy.IsZero() {
			arriveBy := t.LeaveAt
			for i := range route.Legs {
				t.TotalDistance += route.Legs[i].Distance.Meters
				arriveBy = arriveBy.Add(route.Legs[i].Duration * time.Nanosecond)
			}
			t.ArriveBy = arriveBy
		}

		var previousTimeStamp time.Time

		for i, s := range t.Stops {
			if i == 0 {
				s.TimeStamp = t.LeaveAt
			} else if i == (len(t.Stops) - 1) {
				s.TimeStamp = t.ArriveBy
			} else {
				s.TimeStamp = previousTimeStamp.Add(route.Legs[i].Duration * time.Nanosecond)
			}

			previousTimeStamp = s.TimeStamp
		}

		t.TotalTripPrice = float64(t.TotalDistance/1000.0) * PricePerKilometer
	} else {
		return fmt.Errorf("trip.GoogleMapsRepository: no trips found in google map repository")
	}

	return nil
}
