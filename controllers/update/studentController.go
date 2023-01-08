package update

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/SowinskiBraeden/school-management-api/controllers"
	"github.com/SowinskiBraeden/school-management-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

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
				"Personal.firstname":  data["firstname"],
				"Personal.middlename": data["middlename"],
				"Personal.lastname":   data["lastname"],
				"updated_at":          update_time,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"Personal.firstname": data["firstname"],
				"Personal.lastname":  data["lastname"],
				"updated_at":         update_time,
			},
		}
	}

	result, updateErr := StudentCollection.UpdateOne(
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

	_, updateErr := StudentCollection.UpdateOne(
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

	_, updateErr := StudentCollection.UpdateOne(
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
	findErr := StudentCollection.FindOne(ctx, bson.M{"schooldata.sid": claims.Issuer}).Decode(&student)
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
			"Account.password":     student.HashPassword(data["newpassword1"]),
			"Account.temppassword": false, // If it were a temp password, its not now
			"updated_at":           update_time,
		},
		"$push": bson.M{
			"Account.hashhistory": student.HashPassword(data["newpassword1"]),
		},
	}

	_, updateErr := StudentCollection.UpdateOne(
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
	receiver := student.Personal.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/selfPasswordChanged.html", map[string]string{"username": student.Personal.FirstName}); !sent {
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
	findErr := StudentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": data["sid"]}).Decode(&student)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "student not found",
		})
	}

	if student.Personal.Email != data["email"] {
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
			"Account.password":     student.HashPassword(tempPass),
			"Account.temppassword": true,
			"updated_at":           update_time,
		},
	}

	result, updateErr := StudentCollection.UpdateOne(
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
	receiver := student.Personal.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/passwordChanged.html", map[string]string{"username": student.Personal.FirstName, "password": tempPass}); !sent {
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
	err := LockerCollection.FindOne(ctx, bson.M{"lockernumber": data["lockernumber"]}).Decode(&locker)
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

	_, updateErr := StudentCollection.UpdateOne(
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
			"Personal.address":  data["address"],
			"Personal.city":     data["city"],
			"Personal.province": data["province"],
			"Personal.postal":   data["postal"],
			"updated_at":        update_time,
		},
	}

	_, updateErr := StudentCollection.UpdateOne(
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
	findErr := StudentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": data["sid"].(string)}).Decode(&student)
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

	result, updateErr := StudentCollection.UpdateOne(
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
	err := ContactCollection.FindOne(ctx, bson.M{"_id": data["contactid"]}).Decode(&contact)
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
			"Personal.contacts": contact.ID,
		},
	}

	result, updateErr := StudentCollection.UpdateOne(
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
	err := ContactCollection.FindOne(ctx, bson.M{"_id": data["contactid"]}).Decode(&contact)
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
			"Personal.contacts": contact.ID,
		},
	}

	result, updateErr := StudentCollection.UpdateOne(
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
	findErr := StudentCollection.FindOne(context.TODO(), bson.M{"schooldata.sid": sid}).Decode(&student)
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
	findErr = ImageCollection.FindOne(context.TODO(), bson.M{"name": student.SchoolData.PhotoName}).Decode(&photo)
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
	bytes, err := os.ReadFile(fmt.Sprintf("./database/images/%s", image))
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
	_, updateErr := ImageCollection.UpdateOne(
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
			"Personal.email": data["email"],
			"updated_at":     update_time,
		},
	}

	result, updateErr := StudentCollection.UpdateOne(
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
