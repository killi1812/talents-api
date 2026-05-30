package models

import "github.com/google/uuid"

type BookVersion string

type Talent struct {
	Id          uuid.UUID   `bson:"_id,omitempty" json:"id"`
	Name        string      `bson:"name" json:"name"`
	Version     BookVersion `bson:"version" json:"version"`
	MaxLvl      int         `bson:"max_lvl" json:"maxLvl"`
	Test        string      `bson:"test,omitempty" json:"test,omitempty"`
	Description string      `bson:"description" json:"description"`
}
