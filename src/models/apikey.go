package models

import "time"

type APIKey struct {
	Key         string     `bson:"key" json:"key"`
	Description string     `bson:"description" json:"description"`
	LastUsed    *time.Time `bson:"last_used,omitempty" json:"lastUsed,omitempty"`
}
