package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID         primitive.ObjectID `bson:"_id"`
	FirstName  string             `json:"firstname" validate:"required"`
	MiddleName string             `json:"middlename"`
	LastName   string             `json:"lastname" validate:"required"`
	Email      string             `json:"email" validate:"required"`
	Password   string             `json:"password" validate:"required"`
	TID        string             `json:"tid"` // Teacher ID
	Homeroom   string             `json:"homeroom"`
	Address    string             `json:"address"`
	City       string             `json:"city"`
	Province   string             `json:"province"`
	PC         string             `json:"pc"` // Postal Code
	DOB        string             `json:"dob" validate:"required"`
	Photo      string             `json:"string"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
