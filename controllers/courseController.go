package controllers

import (
	"github.com/SowinskiBraeden/school-management-api/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var courseCollection *mongo.Collection = database.OpenCollection(database.Client, "courses")

func CreateCourse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}

func DeleteCourse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}

func UpdateCourseCode(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}

func UpdateCourseName(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}

func UpdateCourseCredit(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}

func UpdateCourseGradelevel(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": true,
	})
}
