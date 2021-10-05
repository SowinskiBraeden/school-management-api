package routes

import (
	"school-management/controllers"

	"github.com/gofiber/fiber/v2"
)

type Person struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

func Setup(app *fiber.App) {
	app.Post("/api/enroll", controllers.Enroll)
}
