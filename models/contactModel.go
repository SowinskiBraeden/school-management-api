package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID         primitive.ObjectID `bson:"_id"`
	FirstName  string             `json:"firstname" validate:"required"`
	MiddleName string             `json:"middlename"`
	LastName   string             `json:"lastname" validate:"required"`
	Province   string             `json:"province"`
	City       string             `json:"city"`
	Address    string             `json:"address" validate:"required"`
	Postal     string             `json:"postal"`
	HomePhone  float64            `json:"homephone" validate:"required"`
	WorkPhone  float64            `json:"workphone"`
	Email      string             `json:"email" validate:"required"`
	Relation   string             `json:"relation" validate:"required"`
	Priotrity  float64            `json:"priority" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
