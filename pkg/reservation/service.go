package reservation

import (
	"fmt"

	"azure.com/ecovo/trip-service/pkg/entity"
	"azure.com/ecovo/trip-service/pkg/trip"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves reservations.
type UseCase interface {
	Register(r *entity.Reservation) error
	Delete(r *entity.Reservation) error
}

// A Service handles the business logic related to reservations.
type Service struct {
	tripService trip.UseCase
}

// NewService creates a reservation service to handle business logic and manipulate
// reservations through a repository.
func NewService(tripService trip.UseCase) *Service {
	return &Service{tripService}
}

// Register modifies trip repository based on a reservation done.
func (s *Service) Register(r *entity.Reservation) error {
	if r == nil {
		return fmt.Errorf("reservation.Service: reservation is nil")
	}

	err := r.Validate()
	if err != nil {
		return err
	}

	t, err := s.tripService.FindByID(r.TripID)
	if err != nil {
		return err
	}

	// We remove seats we want to reserve on the trip
	isInTrip := false
	for _, s := range t.Stops {
		if r.SourceID == s.ID {
			isInTrip = true
		} else if r.DestinationID == s.ID {
			isInTrip = false
		}

		if isInTrip {
			if s.Seats < r.Seats {
				return fmt.Errorf("not enough space in the car")
			}
			s.Seats -= r.Seats
		}
	}

	err = s.tripService.Update(t)
	if err != nil {
		return err
	}

	return nil
}

// Delete modifies trip respository based on a reservation done.
func (s *Service) Delete(r *entity.Reservation) error {
	if r == nil {
		return fmt.Errorf("reservation.Service: reservation is nil")
	}

	err := r.Validate()
	if err != nil {
		return err
	}

	t, err := s.tripService.FindByID(r.TripID)
	if err != nil {
		return err
	}

	// We remove seats we want to reserve on the trip
	isInTrip := false
	for _, s := range t.Stops {
		if r.SourceID == s.ID {
			isInTrip = true
		} else if r.DestinationID == s.ID {
			isInTrip = false
		}

		if isInTrip {
			if (s.Seats + r.Seats) > t.Seats {
				return fmt.Errorf("can't add more seats than the car has")
			}
			s.Seats += r.Seats
		}
	}

	err = s.tripService.Update(t)
	if err != nil {
		return err
	}

	return nil
}
