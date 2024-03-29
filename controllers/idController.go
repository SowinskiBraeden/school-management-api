package controllers

import (
	"context"
	"crypto/rand"
	"io"
	"time"

	"github.com/SowinskiBraeden/school-management-api/database"
	"github.com/SowinskiBraeden/school-management-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var IdCollection *mongo.Collection = database.OpenCollection(database.Client, "cids")

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func ValidateID(id string, userType int) bool { // true: valid id, false: id already in use
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	var foundID models.Id
	err := IdCollection.FindOne(ctx, bson.M{"cid": id}).Decode(&foundID)
	cancel()
	if err != nil { // If there is no id found create new ID object to be stored and return true (unless insert error then try again)
		var newID models.Id
		newID.CID = id
		newID.ParentType = userType
		newID.ID = primitive.NewObjectID()
		_, insertErr := IdCollection.InsertOne(context.Background(), newID)
		return insertErr == nil
	}
	return false
}

func ValidatePEN(pen string) bool { // true: valid pen, false: pen already in use
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	var foundID models.Id
	err := StudentCollection.FindOne(ctx, bson.M{"school.pen": pen}).Decode(&foundID)
	cancel()
	return err != nil
}

func GenerateID(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
