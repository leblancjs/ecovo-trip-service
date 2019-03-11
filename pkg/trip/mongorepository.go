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

// A MongoRepository is a repository that performs CRUD operations on trips in
// a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Driver      *entity.Driver     `bson:"driver"`
	Vehicle     *entity.Vehicle    `bson:"vehicle"`
	Source      *entity.Map        `bson:"source"`
	Destination *entity.Map        `bson:"destination"`
	LeaveAt     time.Time          `bson:"leaveAt"`
	ArriveBy    time.Time          `bson:"arriveBy"`
	Seats       int                `bson:"seats"`
	Stops       []*entity.Map      `bson:"stops"`
	Details     *entity.Details    `bson:"details"`
}

func newDocumentFromEntity(t *entity.Trip) (*document, error) {
	if t == nil {
		return nil, fmt.Errorf("trop.MongoRepository: entity is nil")
	}

	var id primitive.ObjectID
	if t.ID.IsZero() {
		id = primitive.NilObjectID
	} else {
		objectID, err := primitive.ObjectIDFromHex(t.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("trip.MongoRepository: failed to create object")
		}

		id = objectID
	}

	var driverID primitive.ObjectID
	driverID, err := primitive.ObjectIDFromHex(t.Driver.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: failed to create object")
	}

	var vehicleID primitive.ObjectID
	vehicleID, err = primitive.ObjectIDFromHex(t.Vehicle.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: failed to create object")
	}

	t.Driver.ID = driverID
	t.Vehicle.ID = vehicleID

	return &document{
		id,
		t.Driver,
		t.Vehicle,
		t.Source,
		t.Destination,
		t.LeaveAt,
		t.ArriveBy,
		t.Seats,
		t.Stops,
		t.Details,
	}, nil
}

func (d document) Entity() *entity.Trip {
	return &entity.Trip{
		entity.NewIDFromHex(d.ID.Hex()),
		d.Driver,
		d.Vehicle,
		d.Source,
		d.Destination,
		d.LeaveAt,
		d.ArriveBy,
		d.Seats,
		d.Stops,
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

// Find retrieves all trips.
func (r *MongoRepository) Find() ([]*entity.Trip, error) {
	findOptions := options.Find()
	filter := bson.D{{}}
	cur, err := r.collection.Find(context.TODO(), filter, findOptions)

	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: no trip found (%s)", err)
	}

	var trips []*entity.Trip
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

// FindByUserID retrieves all trips from a user.
func (r *MongoRepository) FindByUserID(DriverID entity.ID) ([]*entity.Trip, error) {
	objectID, err := primitive.ObjectIDFromHex(string(DriverID))
	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: failed to create object ID")
	}

	findOptions := options.Find()
	filter := bson.D{{"driver.id", objectID}}

	cur, err := r.collection.Find(context.TODO(), filter, findOptions)

	if err != nil {
		return nil, fmt.Errorf("trip.MongoRepository: no trip found (%s)", err)
	}

	var trips []*entity.Trip
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
	objectID, err := primitive.ObjectIDFromHex(string(ID))
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
