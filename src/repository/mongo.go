package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func ConnectMongo(uri string) (*mongo.Client, error) {
	zap.S().Infof("Connecting to MongoDB at %s", uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		zap.S().Errorf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		zap.S().Errorf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	zap.S().Info("Successfully connected to MongoDB")
	return client, nil
}
