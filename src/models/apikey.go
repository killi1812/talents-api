package models

type APIKey struct {
	Key         string `bson:"key" json:"key"`
	Description string `bson:"description" json:"description"`
}
