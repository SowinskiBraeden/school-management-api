package update

import (
	"context"
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

	_, updateErr := TeacherCollection.UpdateOne(
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
	findErr := TeacherCollection.FindOne(ctx, bson.M{"schooldata.tid": claims.Issuer}).Decode(&teacher)
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

	_, updateErr := TeacherCollection.UpdateOne(
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
	receiver := teacher.Personal.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/selfPasswordChanged.html", map[string]string{"username": teacher.Personal.FirstName}); !sent {
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
	findErr := TeacherCollection.FindOne(context.TODO(), bson.M{"schooldata.tid": data["tid"]}).Decode(&teacher)
	if findErr != nil {
		cancel()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "teacher not found",
		})
	}

	if teacher.Personal.Email != data["email"] {
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

	result, updateErr := StudentCollection.UpdateOne(
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
	receiver := teacher.Personal.Email
	r := NewRequest([]string{receiver}, subject)

	if sent := r.Send("./templates/passwordChanged.html", map[string]string{"username": teacher.Personal.FirstName, "password": tempPass}); !sent {
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
			"Personal.address":  data["address"],
			"Personal.city":     data["city"],
			"Personal.province": data["province"],
			"Personal.postal":   data["postal"],
			"updated_at":        update_time,
		},
	}

	_, updateErr := TeacherCollection.UpdateOne(
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
	findErr := TeacherCollection.FindOne(context.TODO(), bson.M{"schooldata.tid": tid}).Decode(&teacher)
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
	findErr = ImageCollection.FindOne(context.TODO(), bson.M{"name": teacher.SchoolData.PhotoName}).Decode(&photo)
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
			"Personal.email": data["email"],
			"updated_at":     update_time,
		},
	}

	_, updateErr := TeacherCollection.UpdateOne(
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
	findErr := TeacherCollection.FindOne(ctx, bson.M{"schooldata.tid": data["tid"]}).Decode(&teacher)
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
			"Personal.firstname":  data["firstname"],
			"Personal.middlename": middlename,
			"Personal.lastname":   data["lastname"],
			"updated_at":          update_time,
		},
	}

	_, updateErr := TeacherCollection.UpdateOne(
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
