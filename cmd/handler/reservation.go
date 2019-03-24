package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/trip-service/pkg/entity"
	"azure.com/ecovo/trip-service/pkg/reservation"
)

// CreateReservation handles a request to create a reservation.
func CreateReservation(service reservation.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var res *entity.Reservation
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			return err
		}

		err = service.Register(res)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		return nil
	}
}

// DeleteReservation handles a request to delete a reservation.
func DeleteReservation(service reservation.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var res *entity.Reservation
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			return err
		}

		err = service.Delete(res)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}
