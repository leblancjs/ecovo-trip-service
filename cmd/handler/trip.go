package handler

import (
	"net/http"

	"azure.com/ecovo/trip-service/pkg/trip"
)

// CreateTrip handles a request to create a trip.
func CreateTrip(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO - Create trip

		return nil
	}
}

// DeleteTrip handles a request to delete a trip by its unique identifier.
func DeleteTrip(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO - Delete trip by its ID

		return nil
	}
}

// GetTripByID handles a request to retrieve a trip by its unique identifier.
func GetTripByID(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO - Get a trip by its ID

		return nil
	}
}

// GetTrips handles a request to retrieve a trip by its unique identifier.
func GetTrips(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// TODO - Get list of trips

		return nil
	}
}
