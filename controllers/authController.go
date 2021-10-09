package controllers

import (
	"context"
	"school-management/database"
	"school-management/models"
	"time"

	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")
var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")

const SecretKey = "secret"

func Enroll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check minimum enroll field requirements are met
	if data["firstname"] == "" || data["lastname"] == "" || data["age"] == "" || data["gradelevel"] == "" || data["dob"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
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

	result, insertErr := studentCollection.InsertOne(ctx, student)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be inserted",
			"error":   insertErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": true,
		"message": "successfully inserted student",
		"result":  result,
	})
}

func RegisterTeacher(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check minimum register teacher field requirements are met
	if data["firstname"] == "" || data["lastname"] == "" || data["dob"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
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

	result, insertErr := teacherCollection.InsertOne(ctx, teacher)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher could not be inserted",
			"error":   insertErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully inserted teacher",
		"result":  result,
	})
}

func StudentLogin(c *fiber.Ctx) error {
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["password"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var student models.Student
	err := studentCollection.FindOne(ctx, bson.M{"sid": data["sid"]}).Decode(&student)
	defer cancel()

	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
			"error":   err,
		})
	}
	defer cancel()

	var verified bool = student.ComparePasswords(data["password"])
	if verified == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    student.SID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 Day
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not log in",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "correct password",
	})
}

func TeacherLogin(c *fiber.Ctx) error {
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check required fields are included
	if data["tid"] == "" || data["password"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var teacher models.Teacher
	err := studentCollection.FindOne(ctx, bson.M{"tid": data["tid"]}).Decode(&teacher)
	defer cancel()

	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
			"error":   err,
		})
	}
	defer cancel()

	var verified bool = teacher.ComparePasswords(data["password"])
	if verified == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    teacher.TID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 Day
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not log in",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "correct password",
	})
}

func Student(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "not authorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"sid": claims.Issuer}).Decode(&student)
	if findErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "successfully logged into student",
		"result":  student,
	})
}

func Teacher(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "not authorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var teacher models.Teacher
	findErr := teacherCollection.FindOne(context.TODO(), bson.M{"tid": claims.Issuer}).Decode(&teacher)
	if findErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "successfully logged into teacher",
		"result":  teacher,
	})
}

func UpdateStudentName(c *fiber.Ctx) error {
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check id and names are included
	if data["_id"] == "" || data["firstname"] == "" || data["middlename"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	studentObjectId, err := primitive.ObjectIDFromHex(data["_id"])
	if err != nil {
		cancel()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
			"error":   err,
		})
	}
	update := bson.M{
		"$set": bson.M{
			"firstname":  data["firstname"],
			"middlename": data["middlename"],
			"lastname":   data["lastname"],
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"_id": studentObjectId},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student",
		"result":  result,
	})
}

func UpdateStudentGradeLevel(c *fiber.Ctx) error {
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check required fields are included
	if data["_id"] == "" || data["gradelevel"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	studentObjectId, err := primitive.ObjectIDFromHex(data["_id"])
	if err != nil {
		cancel()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
			"error":   err,
		})
	}
	update := bson.M{
		"$set": bson.M{
			"gradelevel": data["gradelevel"],
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"_id": studentObjectId},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student",
		"result":  result,
	})
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
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check id and names are included
	if data["_id"] == "" || data["firstname"] == "" || data["middlename"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	teacherObjectId, err := primitive.ObjectIDFromHex(data["_id"])
	if err != nil {
		cancel()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
			"error":   err,
		})
	}
	update := bson.M{
		"$set": bson.M{
			"firstname":  data["firstname"],
			"middlename": data["middlename"],
			"lastname":   data["lastname"],
		},
	}

	result, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"_id": teacherObjectId},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated teacher",
		"result":  result,
	})
}
