package entity

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Vehicle contains a vehicle's information.
type Vehicle struct {
	ID primitive.ObjectID `json:"id"`
}

// Validate validates that the vehicle's required fields are filled out correctly.
func (t *Vehicle) Validate() error {
	return nil
}
