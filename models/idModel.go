package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Id struct {
	ID  primitive.ObjectID `bson:"_id"`
	CID string             `bson:"cid"` // custom id for admin, teahcer or student
}
