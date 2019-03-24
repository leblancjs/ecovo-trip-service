package trip

import (
	"azure.com/ecovo/trip-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on trips in a database.
type Repository interface {
	FindByID(ID entity.ID) (*entity.Trip, error)
	Find(filters *entity.Filters) ([]*entity.Trip, error)
	Create(trip *entity.Trip) (entity.ID, error)
	Update(trip *entity.Trip) error
	Delete(ID entity.ID) error
}
