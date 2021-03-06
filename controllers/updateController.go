package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/SowinskiBraeden/school-management-api/database"
	"github.com/SowinskiBraeden/school-management-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var lockerCollection *mongo.Collection = database.OpenCollection(database.Client, "lockers")
var imageCollection *mongo.Collection = database.OpenCollection(database.Client, "images")

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
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

	// Ensure Authenticated admin sent request
	if verified, _ := AuthenticateUser(c, 3); !verified {
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

	var updateMiddle bool = false
	if data["middlename"] != "" {
		updateMiddle = true
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var update bson.M
	if updateMiddle {
		update = bson.M{
			"$set": bson.M{
				"personaldata.firstname":  data["firstname"],
				"personaldata.middlename": data["middlename"],
				"personaldata.lastname":   data["lastname"],
				"updated_at":              update_time,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"personaldata.firstname": data["firstname"],
				"personaldata.lastname":  data["lastname"],
				"updated_at":             update_time,
			},
		}
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
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
	var data map[string]interface{}
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == nil || data["gradelevel"] == nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"schooldata.gradelevel": data["gradelevel"].(float64),
			"updated_at":            update_time,
		},
	}

	_, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"].(string)},
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
	})
}

/*
	Far later on this function is going to be completely automated.
	Instead of an admin sending a request to update the homeroom of
	a student or teacher, the system will take the room number of
	the teacher's or student's Block 2 class from their schedule.

	Though this function would remain for students only, for example
	a student requests a course change, if its their block 2 the
	admin would have to alter their homeroom to be the new class
	number.
*/
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
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
			"schooldata.homeroom": data["homeroom"],
			"updated_at":          update_time,
		},
	}

	_, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
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
	findErr := studentCollection.FindOne(ctx, bson.M{"schooldata.sid": claims.Issuer}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	// Check required fields are included
	if data["password"] == "" || data["newpassword1"] == "" || data["newpassword2"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	if !student.ComparePasswords(data["password"]) {
		cancel()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Your password is incorrect",
		})
	}

	if data["newpassword1"] != data["newpassword2"] {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new password must match",
		})
	}

	if student.UsedPassword(data["newpassword1"]) {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new password cannot be the same as a previous password",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"accountdata.password":     student.HashPassword(data["newpassword1"]),
			"accountdata.temppassword": false, // If it were a temp password, its not now
			"updated_at":               update_time,
		},
		"$push": bson.M{
			"accountdata.hashhistory": student.HashPassword(data["newpassword1"]),
		},
	}

	_, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": claims.Issuer},
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

	// Alert email the password has changed
	subject := "Password Changed"
	receiver := student.PersonalData.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/selfPasswordChanged.html", map[string]string{"username": student.PersonalData.FirstName}); !sent {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to send email to student",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student password",
	})
}

// This is for students to reset their password if they are unable to login
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

	// Check required fields are included (email must be personal email)
	if data["sid"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": data["sid"]}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	if student.PersonalData.Email != data["email"] {
		cancel()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Your personal email is incorrect",
		})
	}

	tempPass := student.GeneratePassword(12, 1, 1, 1)
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"accountdata.password":     student.HashPassword(tempPass),
			"accountdata.temppassword": true,
			"updated_at":               update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
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
	subject := "Password Changed"
	receiver := student.PersonalData.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/passwordChanged.html", map[string]string{"username": student.PersonalData.FirstName, "password": tempPass}); !sent {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to students email",
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and locker are included
	if data["sid"] == "" || data["lockernumber"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var locker models.Locker
	err := lockerCollection.FindOne(ctx, bson.M{"lockernumber": data["lockernumber"]}).Decode(&locker)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "locker not found",
			"error":   err,
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"schooldata.locker": locker.ID,
			"updated_at":        update_time,
		},
	}

	_, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
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
	})
}

func UpdateStudentAddress(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["address"] == "" || data["city"] == "" || data["province"] == "" || data["postal"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"personaldata.address":  data["address"],
			"personaldata.city":     data["city"],
			"personaldata.province": data["province"],
			"personaldata.postal":   data["postal"],
			"updated_at":            update_time,
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
		"message": "successfully updated student",
	})
}

// In the case a student gets held back a grade, we need to update their YOG (Year of Graduation)
func UpdateStudentYOG(c *fiber.Ctx) error {
	var data map[string]interface{}
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["sid"] == "" || data["yog"] == nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": data["sid"].(string)}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"schooldata.yog": data["yog"].(int),
			"updated_at":     update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"].(string)},
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

func RemoveStudentContact(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and contact id are included
	if data["sid"] == "" || data["contactid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var contact models.Contact
	err := contactCollection.FindOne(ctx, bson.M{"_id": data["contactid"]}).Decode(&contact)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "contact not found",
			"error":   err,
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"updated_at": update_time,
		},
		"$pull": bson.M{
			"personaldata.contacts": contact.ID,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the contact could not be added",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully added contact",
		"result":  result,
	})
}

func AddStudentContact(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and contact id are included
	if data["sid"] == "" || data["contactid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var contact models.Contact
	err := contactCollection.FindOne(ctx, bson.M{"_id": data["contactid"]}).Decode(&contact)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "contact not found",
			"error":   err,
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"updated_at": update_time,
		},
		"$push": bson.M{
			"personaldata.contacts": contact.ID,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the contact could not be added",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully added contact",
		"result":  result,
	})
}

func UpdateStudentPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	//Ensure Authenticated admin sent request
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	sid := c.FormValue("sid")
	if sid == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	// Get student
	var student models.Student
	findErr := studentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": sid}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be found",
			"error":   findErr,
		})
	}

	// Collect image
	file, err := c.FormFile("image")
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be retrieved",
			"error":   err,
		})
	}

	// Get student photo
	var photo models.Photo
	findErr = imageCollection.FindOne(context.TODO(), bson.M{"name": student.SchoolData.PhotoName}).Decode(&photo)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student image could not be found",
			"error":   findErr,
		})
	}

	// Save image to local
	uniqueId := uuid.New()
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := strings.Split(file.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", filename, fileExt)
	err = c.SaveFile(file, fmt.Sprintf("./database/images/%s", image))
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be saved",
			"error":   err,
		})
	}

	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(fmt.Sprintf("./database/images/%s", image))
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be read",
			"error":   err,
		})
	}

	var base64Encoding string = toBase64(bytes)

	// Update image name and base64 data
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"base64":     base64Encoding,
			"updated_at": update_time,
		},
	}
	_, updateErr := imageCollection.UpdateOne(
		ctx,
		bson.M{"_id": photo.ID},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	// Remove local image
	os.Remove(fmt.Sprintf("./database/images/%s", image))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated student photo",
	})
}

func UpdateStudentEmail(c *fiber.Ctx) error {
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

	var sid string
	var verifiedStudent bool
	verifiedAdmin, _ := AuthenticateUser(c, 3)
	verifiedStudent, sid = AuthenticateUser(c, 1)
	// Ensure Authenticated admin sent request
	if !verifiedAdmin && !verifiedStudent {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin or teacher can perform this action",
		})
	}

	if verifiedAdmin && data["sid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	} else if verifiedAdmin {
		sid = data["sid"]
	}

	// Check required fields are included
	if data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"personaldata.email": data["email"],
			"updated_at":         update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": sid},
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
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
			"accountdata.accountdisabled": false,
			"accountdata.alerted":         false,
			"accountdata.attempts":        0,
			"updated_at":                  update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.sid": data["sid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student account could not be enabled",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully enabled student account",
		"result":  result,
	})
}

func RemoveTeachersDisabled(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["tid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"accountdata.accountdisabled": false,
			"accountdata.attempts":        0,
			"updated_at":                  update_time,
		},
	}

	result, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": data["tid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher account could not be enabled",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully enabled teacher account",
		"result":  result,
	})
}

/*
	Far later on this function is going to be completely automated.
	Instead of an admin sending a request to update the homeroom of
	a student or teacher, the system will take the room number of
	the teacher's or student's Block 2 class from their schedule.

	Though this function would remain for students only, for example
	a student requests a course change, if its their block 2 the
	admin would have to alter their homeroom to be the new class
	number.
*/
func UpdateTeacherHomeroom(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["tid"] == "" || data["homeroom"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"schooldata.homeroom": data["homeroom"],
			"updated_at":          update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": data["tid"]},
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
	})
}

func UpdateTeacherPassword(c *fiber.Ctx) error {
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

	var teacher models.Teacher
	findErr := teacherCollection.FindOne(ctx, bson.M{"schooldata.tid": claims.Issuer}).Decode(&teacher)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
		})
	}

	// Check required fields are included
	if data["password"] == "" || data["newpassword1"] == "" || data["newpassword2"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	if !teacher.ComparePasswords(data["password"]) {
		cancel()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Your password is incorrect",
		})
	}

	if data["newpassword1"] != data["newpassword2"] {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new passwords must match",
		})
	}

	if teacher.UsedPassword(data["newpassword1"]) {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new password cannot be the same as a previous password",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"accountdata.password":     teacher.HashPassword(data["newpassword1"]),
			"accountdata.temppassword": false, // If it were a temp password, its not now
			"updated_at":               update_time,
		},
		"$push": bson.M{
			"accountdata.hashhistory": teacher.HashPassword(data["newpassword1"]),
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": claims.Issuer},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher password could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	subject := "Password Changed"
	receiver := teacher.PersonalData.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/selfPasswordChanged.html", map[string]string{"username": teacher.PersonalData.FirstName}); !sent {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to teachers email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated teacher password",
	})
}

func ResetTeacherPassword(c *fiber.Ctx) error {
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

	// Check required fields are included (email must be personal email)
	if data["tid"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var teacher models.Teacher
	findErr := teacherCollection.FindOne(context.TODO(), bson.M{"schooldata.tid": data["tid"]}).Decode(&teacher)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
		})
	}

	if teacher.PersonalData.Email != data["email"] {
		cancel()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Your personal email is incorrect",
		})
	}

	tempPass := teacher.GeneratePassword(12, 1, 1, 1)
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"accountdata.password":     teacher.HashPassword(tempPass),
			"accountdata.temppassword": true,
			"updated_at":               update_time,
		},
	}

	result, updateErr := studentCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": data["tid"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher password could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	// Send teacher personal email temp password
	subject := "Password Changed"
	receiver := teacher.PersonalData.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/passwordChanged.html", map[string]string{"username": teacher.PersonalData.FirstName, "password": tempPass}); !sent {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to teachers email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated teacher password",
		"result":  result,
	})
}

func UpdateTeacherAddress(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["tid"] == "" || data["address"] == "" || data["city"] == "" || data["province"] == "" || data["postal"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"personaldata.address":  data["address"],
			"personaldata.city":     data["city"],
			"personaldata.province": data["province"],
			"personaldata.postal":   data["postal"],
			"updated_at":            update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": data["tid"]},
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
	})
}

func UpdateTeacherPhoto(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	//Ensure Authenticated admin sent request
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	tid := c.FormValue("tid")
	if tid == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	// Get teacher
	var teacher models.Teacher
	findErr := teacherCollection.FindOne(context.TODO(), bson.M{"schooldata.tid": tid}).Decode(&teacher)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "the teacher could not be found",
			"error":   findErr,
		})
	}

	// Collect image
	file, err := c.FormFile("image")
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be retrieved",
			"error":   err,
		})
	}

	// Get student photo
	var photo models.Photo
	findErr = imageCollection.FindOne(context.TODO(), bson.M{"name": teacher.SchoolData.PhotoName}).Decode(&photo)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student image could not be found",
			"error":   findErr,
		})
	}

	// Save image to local
	uniqueId := uuid.New()
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := strings.Split(file.Filename, ".")[1]
	image := fmt.Sprintf("%s.%s", filename, fileExt)
	err = c.SaveFile(file, fmt.Sprintf("./database/images/%s", image))
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be saved",
			"error":   err,
		})
	}

	// Read the entire file into a byte slice
	bytes, err := ioutil.ReadFile(fmt.Sprintf("./database/images/%s", image))
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be read",
			"error":   err,
		})
	}

	var base64Encoding string = toBase64(bytes)

	// Update image name and base64 data
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"base64":     base64Encoding,
			"updated_at": update_time,
		},
	}
	_, updateErr := imageCollection.UpdateOne(
		ctx,
		bson.M{"_id": photo.ID},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the image could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	// Remove local image
	os.Remove(fmt.Sprintf("./database/images/%s", image))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated teacher photo",
	})
}

func UpdateTeacherEmail(c *fiber.Ctx) error {
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

	var tid string
	var verifiedTeacher bool

	verifiedAdmin, _ := AuthenticateUser(c, 3)
	verifiedTeacher, tid = AuthenticateUser(c, 2)
	// Ensure Authenticated admin sent request
	if !verifiedAdmin && !verifiedTeacher {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin or teacher can perform this action",
		})
	}

	if verifiedAdmin && data["tid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	} else if verifiedAdmin {
		tid = data["tid"]
	}

	// Check required fields are included
	if data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"personaldata.email": data["email"],
			"updated_at":         update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": tid},
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id and names are included
	if data["tid"] == "" || data["firstname"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	// Get teacher
	var teacher models.Teacher
	findErr := teacherCollection.FindOne(ctx, bson.M{"schooldata.tid": data["tid"]}).Decode(&teacher)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "the teacher could not be found",
			"error":   findErr,
		})
	}

	var middlename string = ""

	if data["middlename"] != "" {
		middlename = data["middlename"]
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"personaldata.firstname":  data["firstname"],
			"personaldata.middlename": middlename,
			"personaldata.lastname":   data["lastname"],
			"updated_at":              update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"schooldata.tid": data["tid"]},
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
	})
}

func UpdateContactName(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["_id"] == "" || data["firstname"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var middlename string = ""

	if data["middlename"] != "" {
		middlename = data["middlename"]
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"firstname":  data["firstname"],
			"middlename": middlename,
			"lastname":   data["lastname"],
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
	})
}

func UpdateContactAddress(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check required fields are included
	if data["_id"] == "" || data["address"] == "" || data["city"] == "" || data["province"] == "" || data["postal"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"address":    data["address"],
			"city":       data["city"],
			"province":   data["province"],
			"postal":     data["postal"],
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
	})
}

func UpdateContactHomePhone(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id of contact and new priority number is included
	if data["_id"] == "" || data["newnumber"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"homephone":  data["newnumber"],
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
	})
}

func UpdateContactWorkPhone(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id of contact and new priority number is included
	if data["_id"] == "" || data["newnumber"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"workphone":  data["newnumber"],
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
	})
}

func UpdateContactEmail(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id of contact and new priority number is included
	if data["_id"] == "" || data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"email":      data["email"],
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
	})
}

func UpdateContactPriority(c *fiber.Ctx) error {
	var data map[string]interface{}
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check id of contact and new priority number is included
	if data["_id"] == nil || data["priority"] == nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	var priority int = data["priority"].(int)

	if priority > 10 || priority < 1 {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid priority",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"priority":   priority,
			"updated_at": update_time,
		},
	}

	_, updateErr := contactCollection.UpdateOne(
		ctx,
		bson.M{"_id": data["_id"]},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "contact could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated contact",
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
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

	_, updateErr := lockerCollection.UpdateOne(
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
	})
}

func RemoveStudent(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check student id is included
	if data["sid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	_, deleteErr := idCollection.DeleteOne(ctx, bson.M{"cid": data["sid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the identification number could not be deleted",
			"error":   deleteErr,
		})
	}

	_, deleteErr = studentCollection.DeleteOne(ctx, bson.M{"schooldata.sid": data["sid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be deleted",
			"error":   deleteErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully deleted student",
	})
}

func RemoveTeacher(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check student id is included
	if data["tid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	_, deleteErr := idCollection.DeleteOne(ctx, bson.M{"cid": data["tid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the identification number could not be deleted",
			"error":   deleteErr,
		})
	}

	_, deleteErr = teacherCollection.DeleteOne(ctx, bson.M{"schooldata.tid": data["tid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the teacher could not be deleted",
			"error":   deleteErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully deleted teacher",
	})
}

func RemoveAdmin(c *fiber.Ctx) error {
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
	if verified, _ := AuthenticateUser(c, 3); !verified {
		cancel()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized: only an admin can perform this action",
		})
	}

	// Check student id is included
	if data["aid"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	_, deleteErr := idCollection.DeleteOne(ctx, bson.M{"cid": data["aid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the identification number could not be deleted",
			"error":   deleteErr,
		})
	}

	_, deleteErr = adminCollection.DeleteOne(ctx, bson.M{"aid": data["aid"]})
	if deleteErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the admin could not be deleted",
			"error":   deleteErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully deleted admin",
	})
}

func UpdateAdminName(c *fiber.Ctx) error {
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

	var admin models.Admin
	findErr := adminCollection.FindOne(ctx, bson.M{"aid": claims.Issuer}).Decode(&admin)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "admin not found",
		})
	}

	// Check required fields are included
	if data["firstname"] == "" || data["lastname"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"firstname":  data["firstname"],
			"lastname":   data["lastname"],
			"updated_at": update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"aid": claims.Issuer},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the admin could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated admin",
	})
}

func UpdateAdminEmail(c *fiber.Ctx) error {
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

	var admin models.Admin
	findErr := adminCollection.FindOne(ctx, bson.M{"aid": claims.Issuer}).Decode(&admin)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "admin not found",
		})
	}

	// Check required fields are included
	if data["email"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"email":      data["email"],
			"updated_at": update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"aid": claims.Issuer},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the admin could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated admin",
	})
}

func UpdateAdminPassword(c *fiber.Ctx) error {
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

	var admin models.Admin
	findErr := adminCollection.FindOne(ctx, bson.M{"aid": claims.Issuer}).Decode(&admin)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "admin not found",
		})
	}

	// Check required fields are included
	if data["password"] == "" || data["newpassword1"] == "" || data["newpassword2"] == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing required fields",
		})
	}

	if !admin.ComparePasswords(data["password"]) {
		cancel()
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Your password is incorrect",
		})
	}

	if data["newpassword1"] != data["newpassword2"] {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Your new passwords must match",
		})
	}

	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"password":     admin.HashPassword(data["newpassword1"]),
			"temppassword": false, // If it were a temp password, its not now
			"updated_at":   update_time,
		},
	}

	_, updateErr := teacherCollection.UpdateOne(
		ctx,
		bson.M{"aid": claims.Issuer},
		update,
	)
	if updateErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the admin password could not be updated",
			"error":   updateErr,
		})
	}
	defer cancel()

	subject := "Password Changed"
	receiver := admin.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/selfPasswordChanged.html", map[string]string{"username": admin.FirstName}); !sent {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Could not send password to admins email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "successfully updated admin password",
	})
}
