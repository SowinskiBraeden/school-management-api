package models

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/SowinskiBraeden/school-management-api/database"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var AdminCollection *mongo.Collection = database.OpenCollection(database.Client, "admins")

type Admin struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"firstname" validate:"required"`
	LastName     string             `json:"lastname" validate:"required"`
	Email        string             `json:"email" validate:"required"`
	SchoolEmail  string             `json:"schoolemail"`
	Password     string             `json:"-" validate:"min=10,max=32"`
	TempPassword bool               `json:"temppassword"`
	AID          string             `json:"aid"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}

func (a *Admin) EmailExists(email string) bool {
	var admin Admin
	findErr := AdminCollection.FindOne(context.TODO(), bson.M{"account.schoolemail": email}).Decode(&admin)
	return findErr == nil
}

func (a *Admin) GenerateSchoolEmail(offset int, lastEmail string) string {
	addr := os.Getenv("SYSTEM_EMAIL_ADDRESS")
	var email string = strings.ToLower(a.LastName) + "_" + strings.ToLower(string(a.FirstName[0])) + addr
	if offset > 0 && offset < len([]rune(a.FirstName))-1 {
		email = lastEmail[:strings.Index(lastEmail, "_")+offset+1] + strings.ToLower(string(a.FirstName[offset])) + addr
	}
	if offset == len([]rune(a.FirstName))-1 {
		email = strings.ToLower(a.LastName) + "_" + strings.ToLower(a.FirstName) + addr
	}
	if offset > len([]rune(a.FirstName))-1 {
		email = strings.ToLower(a.LastName) + "_" + strings.ToLower(a.FirstName) + strconv.Itoa(offset-len([]rune(a.FirstName))) + addr
	}
	return email
}

func (s *Admin) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (a *Admin) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

func (a *Admin) CheckPasswordStrength(password string) bool {

	var hasUpper bool = false
	for _, r := range password {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			hasUpper = true
		}
	}

	var hasLower bool = false
	for _, r := range password {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			hasLower = true
		}
	}

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper && len(password) >= 8 {
		return true
	} else {
		return false
	}
}

func (a *Admin) GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
