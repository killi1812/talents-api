package repository

import (
	"context"

	"talents-api/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *TalentRepository) GetAll(ctx context.Context, limit, skip int64) ([]models.Talent, int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var talents []models.Talent
	if err := cursor.All(ctx, &talents); err != nil {
		return nil, 0, err
	}
	return talents, count, nil
}

func (r *TalentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Talent, error) {
	var talent models.Talent
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&talent)
	if err != nil {
		return nil, err
	}
	return &talent, nil
}

func (r *TalentRepository) Search(ctx context.Context, query string, limit, skip int64) ([]models.Talent, int64, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var talents []models.Talent
	if err := cursor.All(ctx, &talents); err != nil {
		return nil, 0, err
	}
	return talents, count, nil
}

func (r *TalentRepository) Update(ctx context.Context, talent *models.Talent) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": talent.Id}, talent)
	return err
}

func (r *TalentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *TalentRepository) Seed(ctx context.Context, talents []models.Talent) (int, error) {
	var count int
	for _, t := range talents {
		// Check if talent with same name exists
		exists, _ := r.collection.CountDocuments(ctx, bson.M{"name": t.Name})
		if exists == 0 {
			if t.Id == uuid.Nil {
				t.Id = uuid.New()
			}
			_, err := r.collection.InsertOne(ctx, t)
			if err == nil {
				count++
			}
		}
	}
	return count, nil
}
