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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*?,.`~"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

var TeacherCollection *mongo.Collection = database.OpenCollection(database.Client, "teachers")

type Teacher struct {
	ID           primitive.ObjectID `bson:"_id"`
	PersonalData struct {
		FirstName  string `json:"firstname" validate:"required"`
		MiddleName string `json:"middlename"`
		LastName   string `json:"lastname" validate:"required"`
		Email      string `json:"email" validate:"required"`
		Address    string `json:"address"`
		City       string `json:"city"`
		Province   string `json:"province"`
		Postal     string `json:"postal"`
		DOB        string `json:"dob" validate:"required"`
	} `json:"personaldata"`
	SchoolData struct {
		TID       string `json:"tid"` // Teacher ID
		Homeroom  string `json:"homeroom"`
		PhotoName string `json:"photoname"` // name of photo in db
	} `json:"schooldata"`
	AccountData struct {
		VerifiedEmail   bool     `json:"verifiedemail"`
		SchoolEmail     string   `json:"schoolemail"`
		Password        string   `json:"-" validate:"min=10,max=32"`
		AccountDisabled bool     `bson:"accountdisabled"`
		TempPassword    bool     `json:"temppassword"`
		Attempts        int      `json:"attempts"` // login attempts max 5
		HashHistory     []string `json:"-"`        // List of old hashed passwords (not including auto generated passwords)
	} `json:"accountdata"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (t *Teacher) UsedPassword(password string) bool {
	for _, oldHash := range t.AccountData.HashHistory {
		if oldHash == t.HashPassword(password) {
			return true
		}
	}
	return false
}

func (t *Teacher) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (t *Teacher) EmailExists(email string) bool {
	var teacher Teacher
	findErr := TeacherCollection.FindOne(context.TODO(), bson.M{"accountdata.schoolemail": email}).Decode(&teacher)
	return findErr == nil
}

func (t *Teacher) GenerateSchoolEmail(offset int, lastEmail string) string {
	addr := os.Getenv("SYSTEM_EMAIL_ADDRESS")
	var email string = strings.ToLower(t.PersonalData.LastName) + "_" + strings.ToLower(string(t.PersonalData.FirstName[0])) + addr
	if offset > 0 && offset < len([]rune(t.PersonalData.FirstName))-1 {
		email = lastEmail[:strings.Index(lastEmail, "_")+offset+1] + strings.ToLower(string(t.PersonalData.FirstName[offset])) + addr
	}
	if offset == len([]rune(t.PersonalData.FirstName))-1 {
		email = strings.ToLower(t.PersonalData.LastName) + "_" + strings.ToLower(t.PersonalData.FirstName) + addr
	}
	if offset > len([]rune(t.PersonalData.FirstName))-1 {
		email = strings.ToLower(t.PersonalData.LastName) + "_" + strings.ToLower(t.PersonalData.FirstName) + strconv.Itoa(offset-len([]rune(t.PersonalData.FirstName))) + addr
	}
	return email
}

func (t *Teacher) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.AccountData.Password), []byte(password))
	return err == nil
}

func (t *Teacher) CheckPasswordStrength(password string) bool {

	var hasUpper bool = false
	for _, r := range password {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			hasUpper = true
		}
	}

	var hasLower bool = false
	for _, r := range password {
		if unicode.IsLower(r) && unicode.IsLetter(r) {
			hasLower = true
		}
	}

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper && len(password) >= 8 {
		return true
	} else {
		return false
	}
}

func (t *Teacher) GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
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
