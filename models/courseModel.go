package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	Grade_level int                `json:"grade_level" validate:"required"`
	Credits     int                `json:"credits" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
