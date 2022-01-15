package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Locker struct {
	ID             primitive.ObjectID `bson:"_id"`
	LockerNumber   string             `json:"lockernumber"` // Example B123
	LockerPasscode string             `json:"lockerpasscode"`
	LockerType     string             `json:"lockertype"` // Upper / Lower locker
}
