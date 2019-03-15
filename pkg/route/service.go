package route

import (
	"azure.com/ecovo/trip-service/pkg/entity"
)

// UseCase interface
type UseCase interface {
	CreateRoute(t *entity.Trip) error
}

// Service structure
type Service struct {
	repo Repository
}

// NewService creates the service
func NewService(repo Repository) UseCase {
	return &Service{repo}
}

// CreateRoute generates route for a trip
func (s *Service) CreateRoute(t *entity.Trip) error {
	return s.repo.GenerateRoute(t)
}
