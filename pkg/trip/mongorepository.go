package trip

import (
	"context"
	"fmt"
	"time"

	"azure.com/ecovo/trip-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

const (
	// DefaultRadius represents the default radius for a location search.
	DefaultRadius = 2000

	// TimeThreshold represents the time threshold for leaveAt or arriveBy (in hours)
	TimeThreshold = 12
)

// A MongoRepository is a repository that performs CRUD operations on trips in
// a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	DriverID  primitive.ObjectID `bson:"driverId"`
	VehicleID primitive.ObjectID `bson:"vehicleId"`
	Full      bool               `bson:"full"`
	LeaveAt   time.Time          `bson:"leaveAt"`
	ArriveBy  time.Time          `bson:"arriveBy"`
	Seats     int                `bson:"seats"`
	Stops     []*stop            `bson:"stops"`
	Details   *entity.Details    `bson:"details"`
}

type stop struct {
	ID        primitive.ObjectID `bson:"id"`
	Point     *entity.Point      `bson:"point"`
	Seats     int                `bson:"seats"`
	TimeStamp time.Time          `bson:"timestamp"`
}

func newDocumentFromEntity(t *entity.Trip) (*document, error) {
	if t == nil {
		return nil, fmt.Errorf("trop.MongoRepository: entity is nil")
	}

	tripID, err := getObjectID(t.ID)
	if err != nil {
		return nil, err
	}

	driverID, err := getObjectID(t.DriverID)
	if err != nil {
		return nil, err
	}

	vehicleID, err := getObjectID(t.VehicleID)
	if err != nil {
		return nil, err
	}

	stops := make([]*stop, len(t.Stops))
	for i, s := range t.Stops {
		stopID, err := getObjectID(s.ID)
		if err != nil {
			return nil, err
		}

		stops[i] = &stop{
			stopID,
			s.Point,
			s.Seats,
			s.TimeStamp,
		}
	}

	return &document{
		tripID,
		driverID,
		vehicleID,
		t.Full,
		t.LeaveAt,
		t.ArriveBy,
		t.Seats,
		stops,
		t.Details,
	}, nil
}

func (d document) Entity() *entity.Trip {
	stops := make([]*entity.Stop, len(d.Stops))
	for i, s := range d.Stops {
		stops[i] = &entity.Stop{
			entity.NewIDFromHex(s.ID.Hex()),
			s.Point,
			s.Seats,
			s.TimeStamp,
		}
	}

	return &entity.Trip{
		entity.NewIDFromHex(d.ID.Hex()),
		entity.NewIDFromHex(d.DriverID.Hex()),
		entity.NewIDFromHex(d.VehicleID.Hex()),
		d.Full,
		d.LeaveAt,
		d.ArriveBy,
		d.Seats,
		stops,
		d.Details,
	}
}

// NewMongoRepository creates a trip repository for a MongoDB collection.
func NewMongoRepository(collection *mongo.Collection) (Repository, error) {
	if collection == nil {
		return nil, fmt.Errorf("trip.MongoRepository: collection is nil")
	}

	return &MongoRepository{collection}, nil
}

// FindByID retrieves the trip with the given ID, if it exists.
func (r *MongoRepository) FindByID(ID entity.ID) (*entity.Trip, error) {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	var d document
	err = r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: no trip found with ID \"%s\" (%s)", ID, err)
	}

	return d.Entity(), nil
}

// Find retrieves all trips based on given filters.
func (r *MongoRepository) Find(f *entity.Filters) ([]*entity.Trip, error) {
	findOptions := options.Find()

	filter, _ := newDocumentFromFilters(f)

	cur, err := r.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: no trip found (%s)", err)
	}

	trips := make([]*entity.Trip, 0)
	for cur.Next(context.TODO()) {
		var d document
		err := cur.Decode(&d)
		if err != nil {
			return nil, err
		}
		trips = append(trips, d.Entity())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	return trips, nil
}

// Create stores the new trip in the database and returns the unique
// identifier that was generated for it.
func (r *MongoRepository) Create(t *entity.Trip) (entity.ID, error) {
	if t == nil {
		return entity.NilID, fmt.Errorf("trip.MongoRepository: failed to create trip (trip is nil)")
	}

	// Initialising stops data
	for _, s := range t.Stops {
		s.ID = entity.ID(primitive.NewObjectID().Hex())
		s.Seats = t.Seats
	}

	d, err := newDocumentFromEntity(t)
	if err != nil {
		return entity.NilID, fmt.Errorf("trip.MongoRepository: failed to create trip document from entity (%s)", err)
	}

	res, err := r.collection.InsertOne(context.TODO(), d)
	if err != nil {
		return entity.NilID, fmt.Errorf("trip.MongoRepository: failed to create trip (%s)", err)
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.NilID, fmt.Errorf("trip.MongoRepository: failed to get ID of created trip (%s)", err)
	}

	return entity.ID(ID.Hex()), nil
}

// Delete removes the trip with the given ID from the database.
func (r *MongoRepository) Delete(ID entity.ID) error {
	objectID, err := primitive.ObjectIDFromHex(ID.Hex())
	if err != nil {
		return fmt.Errorf("trip.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("trip.MongoRepository: failed to delete trip with ID \"%s\" (%s)", ID, err)
	}

	return nil
}

// Creates a document based on filters
func newDocumentFromFilters(f *entity.Filters) (bson.D, error) {
	d := bson.D{}

	d = append(d, bson.E{"full", false})

	if f.DriverID != "" {
		objectID, err := primitive.ObjectIDFromHex(f.DriverID)
		if err != nil {
			return nil, fmt.Errorf("user.MongoRepository: failed to create object ID")
		}

		d = append(d, bson.E{"driverId", objectID})
	}

	// if f.Seats != nil {
	// 	d = append(d, bson.E{
	// 		"seats", bson.M{
	// 			"$gte": *f.Seats,
	// 		},
	// 	})
	// }

	// if f.DetailsAnimals != nil {
	// 	d = append(d, bson.E{"details.animals", *f.DetailsAnimals})
	// }

	// if f.DetailsLuggages != nil {
	// 	d = append(d, bson.E{
	// 		"details.luggages", bson.M{
	// 			"$lte": *f.DetailsLuggages,
	// 		},
	// 	})
	// }

	if !f.LeaveAt.IsZero() {
		d = append(d, bson.E{
			"leaveAt", bson.M{
				"$gt": f.LeaveAt.Add(time.Hour * (-TimeThreshold)),
				"$lt": f.LeaveAt.Add(time.Hour * TimeThreshold),
			},
		})
	}

	radiusThresh := 0
	if f.RadiusThresh != nil {
		radiusThresh = *f.RadiusThresh
	} else {
		radiusThresh = DefaultRadius
	}

	if f.DestinationLatitude != nil && f.DestinationLongitude != nil {
		d = append(d, bson.E{
			"stops.point", bson.M{
				"$near": bson.M{
					"$geometry": bson.M{
						"type":        "Point",
						"coordinates": []float64{*f.DestinationLongitude, *f.DestinationLatitude},
					},
					"$maxDistance": radiusThresh,
				},
			},
		})
	}

	return d, nil
}

// Gets an object ID from an entity of type ID
func getObjectID(rawID entity.ID) (primitive.ObjectID, error) {
	if rawID.IsZero() {
		return primitive.NilObjectID, nil
	}

	objectID, err := primitive.ObjectIDFromHex(rawID.Hex())
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("trip.MongoRepository: failed to create object")
	}
	return objectID, nil
}
