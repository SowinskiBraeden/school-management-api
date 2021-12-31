package controllers

import (
	"context"
	"os"
	"school-management/database"
	"school-management/models"
	"time"

	"net/smtp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")
var studentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")
var contactCollection *mongo.Collection = database.OpenCollection(database.Client, "contacts")
var adminCollection *mongo.Collection = database.OpenCollection(database.Client, "admins")

const SecretKey = "secret"

var systemEmail string = os.Getenv("SYSTEM_EMAIL")
var systemPassword string = os.Getenv("SYSTEM_PASSWORD")

func Enroll(c *fiber.Ctx) error {
	var data map[string]string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	if err := c.BodyParser(&data); err != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	// Check if admin sent request
	if data["aid"] == "" {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "admin id required",
		})
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"aid": data["aid"]}).Decode(&admin)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid admin id",
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

	var student models.Student
	student.FirstName = data["firstname"]
	student.MiddleName = data["middlename"]
	student.LastName = data["lastname"]
	student.Age, _ = strconv.Atoi(data["age"])
	student.GradeLevel, _ = strconv.Atoi(data["gradelevel"])
	student.DOB = data["dob"]
	student.Email = data["email"]
	student.Province = data["province"]
	student.City = data["city"]
	student.Address = data["address"]
	student.Postal = data["postal"]
	student.Contacts = []string{}

	student.YOG = ((12 - student.GradeLevel) + time.Now().Year()) + 1

	student.SchoolEmail = student.GenerateSchoolEmail()

	tempPass := student.GeneratePassword(12, 1, 1, 1)
	student.Password = student.HashPassword(tempPass)
	student.TempPassword = true
	// Send student personal email temp password

	smtpHost := "smpt.gmail.com"
	smtpPort := "587"

	message := []byte("Your temporary password is: " + tempPass)

	auth := smtp.PlainAuth("", systemEmail, systemPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, systemEmail, []string{student.Email}, message)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to students email",
			"error":   err,
		})
	}

	var sid string
	for {
		sid = GenerateID()
		if ValidateID(sid) == true {
			break
		}
	}
	student.SID = sid

	student.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	student.ID = primitive.NewObjectID()

	_, insertErr := studentCollection.InsertOne(ctx, student)
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
	})
}

func RegisterTeacher(c *fiber.Ctx) error {
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

	// Check if admin sent request
	if data["aid"] == "" {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "admin id required",
		})
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"aid": data["aid"]}).Decode(&admin)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid admin id",
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

	var teacher models.Teacher
	teacher.FirstName = data["firstname"]
	teacher.LastName = data["lastname"]
	teacher.Email = data["email"]

	teacher.SchoolEmail = teacher.GenerateSchoolEmail()

	tempPass := teacher.GeneratePassword(12, 1, 1, 1)
	teacher.Password = teacher.HashPassword(tempPass)
	teacher.TempPassword = true
	// Send teacher personal email temp password

	smtpHost := "smpt.gmail.com"
	smtpPort := "587"

	message := []byte("Your temporary password is: " + tempPass)

	auth := smtp.PlainAuth("", systemEmail, systemPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, systemEmail, []string{teacher.Email}, message)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to teacners email",
			"error":   err,
		})
	}

	var tid string
	// For the unlikely event that an ID is already in use this will simply try again till it gets a id not in use
	for {
		tid = GenerateID()
		isValid := ValidateID(tid)
		if isValid == true {
			break
		}
	}
	teacher.TID = tid

	teacher.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	teacher.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	teacher.ID = primitive.NewObjectID()

	_, insertErr := teacherCollection.InsertOne(ctx, teacher)
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
	})
}

func CreateAdmin(c *fiber.Ctx) error {
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

	// Check if admin sent request
	if data["aid"] == "" {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "admin id required",
		})
	}

	var checkadmin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"aid": data["aid"]}).Decode(&checkadmin)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid admin id",
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

	var admin models.Admin
	admin.FirstName = data["firstname"]
	admin.LastName = data["lastname"]
	admin.Email = data["email"]

	tempPass := admin.GeneratePassword(12, 1, 1, 1)
	admin.Password = admin.HashPassword(tempPass)
	admin.TempPassword = true
	// Send admin personal email temp password

	smtpHost := "smpt.gmail.com"
	smtpPort := "587"

	message := []byte("Your temporary password is: " + tempPass)

	auth := smtp.PlainAuth("", systemEmail, systemPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, systemEmail, []string{admin.Email}, message)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to admins email",
			"error":   err,
		})
	}

	var aid string
	for {
		aid = GenerateID()
		if ValidateID(aid) == true {
			break
		}
	}
	admin.AID = aid

	admin.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.ID = primitive.NewObjectID()

	_, insertErr := adminCollection.InsertOne(ctx, admin)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the admin could not be inserted",
			"error":   insertErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully inserted admin",
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
	err := teacherCollection.FindOne(ctx, bson.M{"tid": data["tid"]}).Decode(&teacher)
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

func AdminLogin(c *fiber.Ctx) error {
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
	if data["aid"] == "" || data["password"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"aid": data["aid"]}).Decode(&admin)
	defer cancel()

	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "admin not found",
			"error":   err,
		})
	}
	defer cancel()

	var verified bool = admin.ComparePasswords(data["password"])
	if verified == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    admin.AID,
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

func Admin(c *fiber.Ctx) error {
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

	var admin models.Admin
	findErr := adminCollection.FindOne(context.TODO(), bson.M{"aid": claims.Issuer}).Decode(&admin)
	if findErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "admin not found",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "successfully logged into admin",
		"result":  admin,
	})
}

// Should work for both teacher and student ends
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully logged out",
	})
}

func CreateContact(c *fiber.Ctx) error {
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

	// Check if admin sent request
	if data["aid"] == "" {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "admin id required",
		})
	}

	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"aid": data["aid"]}).Decode(&admin)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid admin id",
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["firstname"] == "" || data["lastname"] == "" || data["homephone"] == "" || data["email"] == "" || data["priority"] == "" || data["relation"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var contact models.Contact
	contact.FirstName = data["firstname"]
	contact.MiddleName = data["middlename"]
	contact.LastName = data["lastname"]
	contact.HomePhone = data["homephone"]
	contact.WorkPhone = data["workphone"]
	contact.Email = data["email"]
	contact.Province = data["province"]
	contact.City = data["city"]
	contact.Address = data["address"]
	contact.Postal = data["postal"]
	contact.Relation = data["relation"]
	contact.Priotrity, _ = strconv.Atoi(data["priority"])

	contact.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	contact.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	contact.ID = primitive.NewObjectID()

	_, insertErr := contactCollection.InsertOne(ctx, contact)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not insert contact",
			"error":   insertErr,
		})
	}

	update := bson.M{
		"$push": bson.M{
			"contacts": contact.ID,
		},
	}
	_, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"sid": data["sid"]},
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
		"message": "successfully inserted contact to student",
	})
}

func RemoveContact(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}
