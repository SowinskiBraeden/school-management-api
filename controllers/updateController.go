package controllers

import (
	"context"
	"net/smtp"
	"school-management/database"
	"school-management/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var lockerCollection *mongo.Collection = database.OpenCollection(database.Client, "lockers")

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

	// Ensure Authenticated admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and names are included
	// Middle name is optional
	if data["sid"] == "" || data["firstname"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"PersonalData.FirstName":  data["firstname"],
			"PersonalData.middlename": data["middlename"],
			"PersonalData.lastname":   data["lastname"],
			"updated_at":              update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
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

	// Ensure Authorized admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["gradelevel"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"SchoolData.gradelevel": data["gradelevel"],
			"updated_at":            update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
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
		"message": "successfully updated student",
		"result":  result,
	})
}

func UpdateStudentHomeroom(c *fiber.Ctx) error {
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

	// Ensure Authenticated admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["homeroom"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"SchoolData.homeroom": data["homeroom"],
			"updated_at":          update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
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
		"message": "successfully updated student",
		"result":  result,
	})
}

func UpdateStudentPassword(c *fiber.Ctx) error {
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

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "not authorized",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"sid": claims.Issuer}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	// Check required fields are included
	if data["currentPassword"] == "" || data["newPassword1"] == "" || data["newPassword2"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	if student.ComparePasswords(data["currentPassword"]) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Your current password is incorrect",
		})
	}

	if data["newPassword1"] != data["newPassword2"] {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new password must match",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"AccountData.password":     student.HashPassword(data["newPassword1"]),
			"AccountData.temppassword": false, // If it were a temp password, its not now
			"updated_at":               update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student password could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student password",
		"result":  result,
	})
}

func ResetStudentPassword(c *fiber.Ctx) error {
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
	if data["sid"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"sid": data["sid"]}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	if student.PersonalData.Email == data["email"] {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Your personal email is incorrect",
		})
	}

	tempPass := student.GeneratePassword(12, 1, 1, 1)
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"AccountData.password":     student.HashPassword(tempPass),
			"AccountData.temppassword": true,
			"updated_at":               update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student password could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	// Send student personal email temp password
	smtpHost := "smpt.gmail.com"
	smtpPort := "587"

	message := []byte("Your temporary password is: " + tempPass)

	auth := smtp.PlainAuth("", systemEmail, systemPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, systemEmail, []string{student.PersonalData.Email}, message)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to students email",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student password",
		"result":  result,
	})
}

func UpdateStudentLocker(c *fiber.Ctx) error {
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

	// Ensure Authenticated admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and locker are included
	if data["sid"] == "" || data["lockerNumber"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var locker models.Locker
	err := lockerCollection.FindOne(ctx, bson.M{"sid": data["sid"]}).Decode(&locker)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "locker not found",
			"error":   err,
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"School.Data.FirstName": locker.ID,
			"updated_at":            update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
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
		"message": "successfully updated student",
		"result":  result,
	})
}

func UpdateStudentAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateStudentYOG(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateStudentContacts(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateStudentPhoto(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateStudentEmail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func RemoveStudentsDisabled(c *fiber.Ctx) error {
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

	// Ensure Authorized admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"AccountData.accountdisabled": false,
			"AccountData.attempts":        0,
			"updated_at":                  update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student account could not be re-enabled",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully re-enabled student account",
		"result":  result,
	})
}

func UpdateTeacherHomeroom(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateTeacherPassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateTeacherAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateTeacherPhoto(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateTeacherEmail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
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

	// Ensure Authenticated admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
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

	teacherObjectId, idErr := primitive.ObjectIDFromHex(data["_id"])
	if idErr != nil {
		cancel()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
			"error":   idErr,
		})
	}
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"firstname":  data["firstname"],
			"middlename": data["middlename"],
			"lastname":   data["lastname"],
			"updated_at": update_time,
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

func UpdateContactName(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateContactAddress(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateContactPhone(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateContactEmail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateContactPriority(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": nil,
		"message": "not implimented",
	})
}

func UpdateLockerCombo(c *fiber.Ctx) error {
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

	// Ensure Authenticated admin sent request
	if !AuthAdmin(c) {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check locker number is included
	if data["lockernumber"] == "" || data["newlockercombo"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"lockercombo": data["newlockercombo"],
			"updated_at":  update_time,
		},
	}

	result, updateErr := lockerCollection.UpdateOne(
		ctx,
		bson.M{"lockernumber": data["lockernumber"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the locker could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated locker",
		"result":  result,
	})
}