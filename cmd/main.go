package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/trip-service/cmd/handler"
	"azure.com/ecovo/trip-service/cmd/middleware/auth"
	"azure.com/ecovo/trip-service/pkg/db"
	"azure.com/ecovo/trip-service/pkg/trip"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authConfig := auth.Config{
		Domain: os.Getenv("AUTH_DOMAIN")}
	authValidator, err := auth.NewTokenValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}

	dbConnectionTimeout, err := time.ParseDuration(os.Getenv("DB_CONNECTION_TIMEOUT") + "s")
	if err != nil {
		dbConnectionTimeout = db.DefaultConnectionTimeout
	}
	dbConfig := db.Config{
		Host:              os.Getenv("DB_HOST"),
		Username:          os.Getenv("DB_USERNAME"),
		Password:          os.Getenv("DB_PASSWORD"),
		Name:              os.Getenv("DB_NAME"),
		ConnectionTimeout: dbConnectionTimeout}
	db, err := db.New(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	tripRepository, err := trip.NewMongoRepository(db.Trips)
	if err != nil {
		log.Fatal(err)
	}
	tripUseCase := trip.NewService(tripRepository)

	r := mux.NewRouter()

	// Trips
	r.Handle("/trips", handler.RequestID(handler.Auth(authValidator, handler.GetTrips(tripUseCase)))).
		Methods("GET")
	r.Handle("/trips/{id}", handler.RequestID(handler.Auth(authValidator, handler.GetTripByID(tripUseCase)))).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.Handle("/trips", handler.RequestID(handler.Auth(authValidator, handler.CreateTrip(tripUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/trips/{id}", handler.RequestID(handler.Auth(authValidator, handler.DeleteTrip(tripUseCase)))).
		Methods("DELETE").
		HeadersRegexp("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
