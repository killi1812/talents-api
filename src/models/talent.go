package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BookVersion string

type Talent struct {
	Id          uuid.UUID   `bson:"_id,omitempty" json:"id" swaggertype:"string"`
	Name        string      `bson:"name" json:"name"`
	Version     BookVersion `bson:"version" json:"version"`
	MaxLvl      int         `bson:"max_lvl" json:"maxLvl"`
	Test        string      `bson:"test,omitempty" json:"test,omitempty"`
	Description string      `bson:"description" json:"description"`
}

// UnmarshalJSON handles "id": "string" or empty string from Swagger
func (t *Talent) UnmarshalJSON(data []byte) error {
	type Alias Talent
	aux := &struct {
		Id string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Id == "" || aux.Id == "string" {
		t.Id = uuid.Nil
	} else {
		id, err := uuid.Parse(aux.Id)
		if err != nil {
			return err
		}
		t.Id = id
	}
	return nil
}
