package repository

import (
	"context"
	"talents-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type APIKeyRepository struct {
	collection *mongo.Collection
}

func NewAPIKeyRepository(db *mongo.Database) *APIKeyRepository {
	return &APIKeyRepository{
		collection: db.Collection("api_keys"),
	}
}

func (r *APIKeyRepository) Create(ctx context.Context, key *models.APIKey) error {
	_, err := r.collection.InsertOne(ctx, key)
	return err
}

func (r *APIKeyRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"key": key})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
