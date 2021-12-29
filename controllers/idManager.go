package controllers

import (
	"context"
	"crypto/rand"
	"io"
	"school-management/database"
	"school-management/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var idCollection *mongo.Collection = database.OpenCollection(database.Client, "cids")

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func ValidateID(id string) bool { // true: valid id, false: id already in use
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	var foundID models.Id
	err := idCollection.FindOne(ctx, bson.M{"cid": id}).Decode(&foundID)
	cancel()
	if err != nil { // If there is no id found create new ID object to be stored and return true (unless insert error then try again)
		var newID models.Id
		newID.CID = id
		newID.ID = primitive.NewObjectID()
		_, insertErr := idCollection.InsertOne(context.Background(), newID)
		if insertErr != nil {
			return false
		}
		return true
	}
	return false
}

func GenerateID() string {
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if n != 6 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}