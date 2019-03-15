package route

import (
	"azure.com/ecovo/trip-service/pkg/entity"
)

// Repository interface
type Repository interface {
	GenerateRoute(t *entity.Trip) error
}
