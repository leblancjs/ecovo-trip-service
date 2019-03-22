package trip

import (
	"fmt"

	"azure.com/ecovo/trip-service/pkg/entity"
	"azure.com/ecovo/trip-service/pkg/route"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves trips.
type UseCase interface {
	Register(t *entity.Trip) (*entity.Trip, error)
	FindByID(ID entity.ID) (*entity.Trip, error)
	Find(filters *entity.Filters) ([]*entity.Trip, error)
	Delete(ID entity.ID) error
}

// A Service handles the business logic related to trips.
type Service struct {
	repo         Repository
	routeService route.UseCase
}

// NewService creates a trip service to handle business logic and manipulate
// trips through a repository.
func NewService(repo Repository, routeService route.UseCase) *Service {
	return &Service{repo, routeService}
}

// Register validates the trips's information
func (s *Service) Register(t *entity.Trip) (*entity.Trip, error) {
	if t == nil {
		return nil, fmt.Errorf("trip.Service: trip is nil")
	}

	err := t.Validate()
	if err != nil {
		return nil, err
	}

	err = s.routeService.CreateRoute(t)
	if err != nil {
		return nil, err
	}

	t.ID, err = s.repo.Create(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// FindByID retrieves the trip with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.Trip, error) {
	t, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return t, nil
}

// Find retrieves all the trips
func (s *Service) Find(filters *entity.Filters) ([]*entity.Trip, error) {
	err := filters.Validate()
	if err != nil {
		return nil, err
	}

	t, err := s.repo.Find(filters)
	if err != nil {
		return []*entity.Trip{}, err
	}

	return t, nil
}

// Delete erases the trip from the repository.
func (s *Service) Delete(ID entity.ID) error {
	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
