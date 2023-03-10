package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/trip-service/cmd/handler"
	"azure.com/ecovo/trip-service/cmd/middleware/auth"
	"azure.com/ecovo/trip-service/pkg/db"
	"azure.com/ecovo/trip-service/pkg/pubsub"
	"azure.com/ecovo/trip-service/pkg/pubsub/subscription"
	"azure.com/ecovo/trip-service/pkg/reservation"
	"azure.com/ecovo/trip-service/pkg/route"
	"azure.com/ecovo/trip-service/pkg/trip"
	"github.com/ably/ably-go/ably"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"googlemaps.github.io/maps"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authConfig := auth.Config{
		Domain:               os.Getenv("AUTH_DOMAIN"),
		BasicAuthCredentials: os.Getenv("AUTH_CREDENTIALS"),
	}
	authBasicValidator, err := auth.NewBasicAuthValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}
	authTokenValidator, err := auth.NewTokenValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}
	authValidators := map[string]auth.Validator{
		"basic":  authBasicValidator,
		"bearer": authTokenValidator,
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

	ablyClient, err := ably.NewRestClient(ably.NewClientOptions(os.Getenv("ABLY_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	ablyPubSubRepository, err := subscription.NewAblyRepository(ablyClient)
	if err != nil {
		log.Fatal(err)
	}
	pubSubService := pubsub.NewService(ablyPubSubRepository)

	mapsClient, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	routeRepository, err := route.NewGoogleMapsRepository(mapsClient)
	if err != nil {
		log.Fatal(err)
	}
	routeUseCase := route.NewService(routeRepository)

	tripRepository, err := trip.NewMongoRepository(db.Trips)
	if err != nil {
		log.Fatal(err)
	}
	tripUseCase := trip.NewService(tripRepository, pubSubService, routeUseCase)

	reservationUseCase := reservation.NewService(tripUseCase)

	r := mux.NewRouter()

	// Trips
	r.Handle("/trips", handler.RequestID(handler.Auth(authValidators, handler.GetTrips(tripUseCase)))).
		Methods("GET")
	r.Handle("/trips/{id}", handler.RequestID(handler.Auth(authValidators, handler.GetTripByID(tripUseCase)))).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.Handle("/trips", handler.RequestID(handler.Auth(authValidators, handler.CreateTrip(tripUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/trips/{id}/reservation", handler.RequestID(handler.Auth(authValidators, handler.CreateReservation(reservationUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/trips/{id}/reservation", handler.RequestID(handler.Auth(authValidators, handler.DeleteReservation(reservationUseCase)))).
		Methods("DELETE").
		HeadersRegexp("Content-Type", "application/json")
	r.Handle("/trips/{id}", handler.RequestID(handler.Auth(authValidators, handler.DeleteTrip(tripUseCase)))).
		Methods("DELETE").
		HeadersRegexp("Content-Type", "application/json")
	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
