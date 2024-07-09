package services

import (
	"context"

	"c2c.in/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UnitService struct holds the MongoDB database connection
type UnitService struct {
	DB             *mongo.Database
	collectionName string
}

// NewUnitService creates a new UnitService instance
func NewUnitService(db *mongo.Database) *UnitService {
	return &UnitService{DB: db, collectionName: "units"}
}

// CreateUnit method inserts a new unit into the MongoDB collection
// It returns the ID of the inserted unit and an error if there was a problem
func (s *UnitService) CreatUnit(unit *models.Unit) (string, error) {
	collection := s.DB.Collection(s.collectionName)
	result, err := collection.InsertOne(context.Background(), unit)
	if err != nil {
		return "unable to insert the record", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
