package main

import (
	"fmt"
	"os"

	"github.com/SowinskiBraeden/school-management-api/routes"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const version string = "\nv0.8.7-Alpha"

func main() {
	fmt.Println(version)
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
