package services

import (
	"context"

	"c2c.in/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PathwayService struct holds the MongoDB database connection
type PathwayService struct {
	DB             *mongo.Database
	collectionName string
}

func NewPathwayService(db *mongo.Database) *PathwayService {
	return &PathwayService{DB: db, collectionName: "pathways"}
}

func (p *PathwayService) CreatePathway(pathway *models.Pathway) (string, error) {
	collection := p.DB.Collection(p.collectionName)
	result, err := collection.InsertOne(context.Background(), &pathway)
	if err != nil {
		return "unable to create the pathway", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil

}
