package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type BookVersion string

const (
	VersionCore         BookVersion = "core"
	VersionWindsOfMagic BookVersion = "windsOfMagic"
	VersionUpInArms     BookVersion = "upInArms"
)

type Skill string

const (
	SkillWeaponSkill    Skill = "Weapon Skill"
	SkillBallisticSkill Skill = "Ballistic Skill"
	SkillStrength       Skill = "Strength"
	SkillToughness      Skill = "Toughness"
	SkillInitiative     Skill = "Initiative"
	SkillAgility        Skill = "Agility"
	SkillDexterity      Skill = "Dexterity"
	SkillIntelligence   Skill = "Intelligence"
	SkillWillpower      Skill = "Willpower"
	SkillFellowship     Skill = "Fellowship"
)

type Talent struct {
	Id          uuid.UUID   `bson:"_id,omitempty" json:"id" swaggertype:"string"`
	Name        string      `bson:"name" json:"name"`
	Version     BookVersion `bson:"version" json:"version" swaggertype:"string"`
	MaxLvl      string      `bson:"max_lvl" json:"maxLvl"`
	Test        string      `bson:"test,omitempty" json:"test,omitempty"`
	Description string      `bson:"description" json:"description"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// UnmarshalJSON handles "id": "string" or empty string from Swagger and validates Talent fields
func (t *Talent) UnmarshalJSON(data []byte) error {
	type Alias Talent
	aux := &struct {
		Id      string `json:"id"`
		Version string `json:"version"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// ID handling
	if aux.Id == "" || aux.Id == "string" {
		t.Id = uuid.Nil
	} else {
		id, err := uuid.Parse(aux.Id)
		if err != nil {
			return err
		}
		t.Id = id
	}

	// Version handling (default "core")
	if aux.Version == "" {
		t.Version = VersionCore
	} else {
		v := BookVersion(aux.Version)
		switch v {
		case VersionCore, VersionWindsOfMagic, VersionUpInArms:
			t.Version = v
		default:
			return fmt.Errorf("invalid version: %s (allowed: core, windsOfMagic, upInArms)", aux.Version)
		}
	}

	return nil
}
