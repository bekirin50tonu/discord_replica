package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	collection mongo.Collection
}

func NewSessionRepository(db *mongo.Database, collectionName string) (*SessionRepository, error) {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		collection: *db.Collection(collectionName),
	}, nil
}
