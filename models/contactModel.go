package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID             primitive.ObjectID `bson:"_id"`
	FirstName      string             `json:"firstname" validate:"required"`
	MiddleName     string             `json:"middlename"`
	LastName       string             `json:"lastname" validate:"required"`
	Address        string             `json:"address" validate:"required"`
	Phone          string             `json:"phone" validate:"required"`
	SecondaryPhone string             `json:"secondaryphone"`
	Email          string             `json:"email" validate:"required"`
	Relation       string             `json:"relation" validate:"required"`
	PrimaryContact bool               `json:"primaraycontact" validate:"required"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
}
