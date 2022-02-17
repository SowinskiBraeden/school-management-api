package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `json:"name" validate:"required"`
	Code         string             `json:"code" validate:"required"` // Short string to identify each class
	GradeLevel   int                `json:"gradelevel" validate:"required"`
	Credits      int                `json:"credits" default:"4"`
	TotalRequest int                `json:"totalRequest" default:"0"` // Amount of requests for this course
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}
