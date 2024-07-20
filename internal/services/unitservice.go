package services

import (
	"context"
	"errors"

	"c2c.in/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
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

	fileExtension, ok := unit.MetaData["file-extension"]
	if !ok {
		return "", errors.New("file extension not found in MetaData")
	}

	var prefix string
	switch fileExtension {
	case "pdf":
		prefix = "PDF_"

	case "hls":
		prefix = "HLS_"

	default:
		return "", errors.New("unsupported file extension")

	}

	collection := s.DB.Collection(s.collectionName)
	result, err := collection.InsertOne(context.Background(), unit)
	if err != nil {
		return "unable to insert the record", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	finalId := prefix + id
	return finalId, nil
}

func (s *UnitService) GetAllUnits() ([]string, error) {
	collection := s.DB.Collection(s.collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var units []models.Unit
	err = cursor.All(context.Background(), &units)
	if err != nil {
		return nil, err
	}

	var unitNames []string
	for _, unit := range units {
		unitNames = append(unitNames, unit.Name)
	}

	return unitNames, nil

}
