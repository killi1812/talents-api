package repository

import (
	"context"
	"talents-api/models"
	"time"

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

func (r *APIKeyRepository) GetAll(ctx context.Context) ([]models.APIKey, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var keys []models.APIKey
	if err := cursor.All(ctx, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, key string) error {
	now := time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"key": key}, bson.M{"$set": bson.M{"last_used": now}})
	return err
}

func (r *APIKeyRepository) Delete(ctx context.Context, key string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"key": key})
	return err
}
