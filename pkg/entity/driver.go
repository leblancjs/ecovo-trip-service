package entity

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Driver contains a drivers's information.
type Driver struct {
	ID primitive.ObjectID `json:"id"`
}

// Validate validates that the drivers's required fields are filled out correctly.
func (t *Driver) Validate() error {
	return nil
}
