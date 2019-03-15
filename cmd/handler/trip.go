package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/trip-service/pkg/entity"
	"azure.com/ecovo/trip-service/pkg/trip"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// CreateTrip handles a request to create a trip.
func CreateTrip(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var t *entity.Trip
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			return err
		}

		t, err = service.Register(t)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			_ = service.Delete(entity.ID(t.ID))

			return err
		}

		return nil
	}
}

// DeleteTrip handles a request to delete a trip by its unique identifier.
func DeleteTrip(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])

		err := service.Delete(id)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetTripByID handles a request to retrieve a trip by its unique identifier.
func GetTripByID(tService trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])
		t, err := tService.FindByID(id)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetTrips handles a request to retrieve a trip by its unique identifier.
func GetTrips(service trip.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var encoder = schema.NewDecoder()
		var f entity.Filters

		err := encoder.Decode(&f, r.URL.Query())
		if err != nil {
			return err
		}

		t, err := service.Find(&f)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			return err
		}

		return nil
	}
}
