package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	This is a stupid Idea but in the cids collection
	is a list of all id's (student id's [sid], teacher
	id's [tid], or admin id's [aid]) this way when
	generating a new id for a new user, only the cids
	collection has to be searched to see the id is
	already in use. Rather than looking in the admin,
	teacher, and student collections to check an id.
*/

type Id struct {
	ID  primitive.ObjectID `bson:"_id"`
	CID string             `bson:"cid"` // custom id for admin, teahcer or student
}
