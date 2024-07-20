package services

import (
	"context"

	"c2c.in/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModuleService struct holds the MongoDB database connection
type ModuleService struct {
	DB             *mongo.Database
	collectionName string
}

// NewModuleService creates a new ModuleService instance
func NewModuleService(db *mongo.Database) *ModuleService {
	return &ModuleService{DB: db, collectionName: "modules"}
}

// It returns the ID of the inserted unit and an error if there was a problem
func (m *ModuleService) CreateModule(module *models.Module) (string, error) {
	collection := m.DB.Collection(m.collectionName)
	result, err := collection.InsertOne(context.Background(), module)
	if err != nil {
		return "unable to insert the module", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (m *ModuleService) GetAllModules() ([]string, error) {
	collection := m.DB.Collection(m.collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var modules []models.Module
	err = cursor.All(context.Background(), &modules)
	if err != nil {
		return nil, err
	}

	var moduleNames []string
	for _, module := range modules {
		moduleNames = append(moduleNames, module.Name)
	}

	return moduleNames, nil
}
