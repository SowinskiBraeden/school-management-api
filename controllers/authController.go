package controllers

import (
	"school-management/models"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Test struct {
	a string
	b int
}

func (t *Test) getA() string {
	return t.a
}

func Enroll(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Get age and grade level and convert to int
	intAge, _ := strconv.Atoi(data["age"])
	intGradeLevel, _ := strconv.Atoi(data["gradelevel"])

	student := models.Student{
		FirstName:  data["firstname"],
		LastName:   data["lastname"],
		Age:        intAge,
		GradeLevel: intGradeLevel,
		DOB:        data["dob"],
	}
	student.Password = student.HashPassword(data["password"])
	student.Email = student.GenerateSchoolEmail()

	return c.JSON(student)
}
