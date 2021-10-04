package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"school-management/models"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load(".env")
	mongoURI := os.Getenv("mongoURI")
	dbo := os.Getenv("dbo")

	// Connect to mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// Get collection
	collection := client.Database(dbo).Collection("students")
	// res, insertErr := collection.InsertOne(ctx, student)
	// if insertErr != nil {
	// 	log.Fatal(insertErr)
	// }
	// fmt.Println(res)

	cur, currErr := collection.Find(ctx, bson.D{})
	if currErr != nil {
		panic(currErr)
	}
	defer cur.Close(ctx)

	var students []models.Student
	if err = cur.All(ctx, &students); err != nil {
		panic(err)
	}
	fmt.Println(students)
}
