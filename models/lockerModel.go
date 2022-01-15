package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Locker struct {
	ID           primitive.ObjectID `bson:"_id"`
	LockerNumber string             `json:"lockernumber"` // Example B123
	LockerCombo  string             `json:"lockercombo"`
	LockerType   string             `json:"lockertype"` // Upper / Lower locker
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}
