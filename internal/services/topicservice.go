package services

import (
	"context"

	"c2c.in/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TopicService struct holds the MongoDB database connection
type TopicService struct {
	DB             *mongo.Database
	collectionName string
}

// NewTopicService creates a TopicService instance
func NewTopicService(db *mongo.Database) *TopicService {
	return &TopicService{DB: db, collectionName: "topics"}
}

// CreateUnit method inserts a new unit into the MongoDB collection
// It returns the ID of the inserted unit and an error if there was a problem
func (t *TopicService) CreatTopic(topic *models.Topic) (string, error) {
	collection := t.DB.Collection(t.collectionName)
	result, err := collection.InsertOne(context.Background(), topic)
	if err != nil {
		return "unable to create a topic", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (t *TopicService) GetAllTopics() ([]models.Topic, error) {
	collection := t.DB.Collection(t.collectionName)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var topics []models.Topic
	err = cursor.All(context.Background(), &topics)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (t *TopicService) GetSpecificTopic(topicID string) (*models.Topic, error) {

	collection := t.DB.Collection(t.collectionName)

	objectID, err := primitive.ObjectIDFromHex(topicID)

	filter := bson.M{"_id": objectID}

	result := collection.FindOne(context.Background(), filter)

	var topic models.Topic
	err = result.Decode(&topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}
