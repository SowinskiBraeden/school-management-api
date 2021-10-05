package main

import (
	"school-management/database"

	"school-management/routes"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client = database.DBinstance()
var studentCollection *mongo.Collection = database.OpenCollection(client, "students")

func main() {
	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8000")
}
