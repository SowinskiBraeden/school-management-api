package controllers

import (
	"school-management/database"
	"school-management/models"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ErrorStruct struct {
	Error string
}

var teacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")
var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")

func Enroll(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check minimum enroll field requirements are met
	if data["firstname"] == "" || data["lastname"] == "" || data["age"] == "" || data["gradelevel"] == "" || data["dob"] == "" || data["email"] == "" {
		return c.Status(400).JSON(ErrorStruct{Error: "missing required fields"})
	}

	// Get age and grade level and convert to int
	intAge, _ := strconv.Atoi(data["age"])
	intGradeLevel, _ := strconv.Atoi(data["gradelevel"])

	student := models.Student{}
	student.FirstName = data["firstname"]
	student.MiddleName = data["middlename"]
	student.LastName = data["lastname"]
	student.Age = intAge
	student.GradeLevel = intGradeLevel
	student.DOB = data["dob"]
	student.Email = data["email"]
	student.Province = data["province"]
	student.City = data["city"]
	student.Address = data["address"]
	student.Postal = data["postal"]

	student.YOG = ((12 - student.GradeLevel) + time.Now().Year()) + 1

	student.SchoolEmail = student.GenerateSchoolEmail()

	tempPass := student.GeneratePassword(12, 1, 1, 1)
	student.Password = student.HashPassword(tempPass)
	student.TempPassword = true
	// Send student personal email temp password

	student.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.ID = primitive.NewObjectID()

	result, insertErr := studentCollection.InsertOne(c.Context(), student)
	if insertErr != nil {
		return c.Status(500).JSON(ErrorStruct{Error: "the student could not be inserted"})
	}

	return c.Status(201).JSON(result)
}

func RegisterTeacher(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Check minimum register teacher field requirements are met
	if data["firstname"] == "" || data["lastname"] == "" || data["dob"] == "" || data["email"] == "" {
		return c.Status(400).JSON(ErrorStruct{Error: "missing required fields"})
	}

	teacher := models.Teacher{}
	teacher.FirstName = data["firstname"]
	teacher.LastName = data["lastname"]
	teacher.Email = data["email"]

	teacher.SchoolEmail = teacher.GenerateSchoolEmail()

	tempPass := teacher.GeneratePassword(12, 1, 1, 1)
	teacher.Password = teacher.HashPassword(tempPass)
	teacher.TempPassword = true
	// Send teacher personal email temp password

	teacher.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	teacher.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	teacher.ID = primitive.NewObjectID()

	result, insertErr := teacherCollection.InsertOne(c.Context(), teacher)
	if insertErr != nil {
		return c.Status(500).JSON(ErrorStruct{Error: "the teacher could not be inserted"})
	}

	return c.Status(201).JSON(result)
}

func UpdateStudentName(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentGradeLevel(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentHomeroom(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentPassword(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentLocker(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentAddress(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentYOG(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentContacts(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentPhoto(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateStudentEmail(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherHomeroom(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherPassword(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherAddress(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherPhoto(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherEmail(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}

func UpdateTeacherName(c *fiber.Ctx) error {
	return c.JSON("status:ok")
}
