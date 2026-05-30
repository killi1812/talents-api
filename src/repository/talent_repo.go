package repository

import (
	"context"
	"talents-api/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TalentRepository struct {
	collection *mongo.Collection
}

func NewTalentRepository(db *mongo.Database) *TalentRepository {
	return &TalentRepository{
		collection: db.Collection("talents"),
	}
}

func (r *TalentRepository) Create(ctx context.Context, talent *models.Talent) error {
	if talent.Id == uuid.Nil {
		talent.Id = uuid.New()
	}
	_, err := r.collection.InsertOne(ctx, talent)
	return err
}

func (r *TalentRepository) GetAll(ctx context.Context) ([]models.Talent, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var talents []models.Talent
	if err := cursor.All(ctx, &talents); err != nil {
		return nil, err
	}
	return talents, nil
}

func (r *TalentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Talent, error) {
	var talent models.Talent
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&talent)
	if err != nil {
		return nil, err
	}
	return &talent, nil
}

func (r *TalentRepository) Search(ctx context.Context, query string) ([]models.Talent, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var talents []models.Talent
	if err := cursor.All(ctx, &talents); err != nil {
		return nil, err
	}
	return talents, nil
}
