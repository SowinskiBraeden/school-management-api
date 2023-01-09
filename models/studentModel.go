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

var LockerCollection *mongo.Collection = database.OpenCollection(database.Client, "lockers")
var StudentCollection *mongo.Collection = database.OpenCollection(database.Client, "students")

type Student struct {
	ID       primitive.ObjectID `bson:"_id"`
	Personal struct {
		FirstName  string   `json:"firstname" validate:"required"`
		MiddleName string   `json:"middlename"`
		LastName   string   `json:"lastname" validate:"required"`
		Age        float64  `json:"age" validate:"required"`
		Email      string   `json:"email" validate:"required"`
		Address    string   `json:"address"`
		City       string   `json:"city"`
		Province   string   `json:"province"`
		Postal     string   `json:"postal"`
		DOB        string   `json:"dob" validate:"required"`
		Contacts   []string `json:"contacts"` // List of contact ID's rather than contact object
	} `json:"personal"`
	School struct {
		GradeLevel float64 `json:"gradelevel" validate:"required"`
		SID        string  `json:"sid"` // Student ID
		PEN        string  `json:"ped"` // Personal Education Number
		Homeroom   string  `json:"homeroom"`
		Locker     string  `json:"-"`         // Locker ID
		YOG        int     `json:"yog"`       // Year of Graduation
		PhotoName  string  `json:"photoname"` // name of photo in db
	} `json:"School"`
	Account struct {
		VerifiedEmail   bool     `json:"verifiedemail"`
		SchoolEmail     string   `json:"schoolemail"`
		Password        string   `json:"-" validate:"min=10,max=32"`
		AccountDisabled bool     `bson:"accountdisabled"`
		Alerted         bool     `bson:"alerted"`
		TempPassword    bool     `json:"temppassword"`
		Attempts        int      `json:"attempts"` // login attempts max 5
		HashHistory     []string `json:"-"`        // List of old hashed passwords (not including auto generated passwords)
	} `json:"Account"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (s *Student) UsedPassword(password string) bool {
	for _, oldHash := range s.Account.HashHistory {
		if oldHash == s.HashPassword(password) {
			return true
		}
	}
	return false
}

func (s *Student) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (s *Student) EmailExists(email string) bool {
	var student Student
	findErr := StudentCollection.FindOne(context.TODO(), bson.M{"account.schoolemail": email}).Decode(&student)
	return findErr == nil
}

func (s *Student) GenerateSchoolEmail(offset int, lastEmail string) string {
	addr := os.Getenv("SYSTEM_EMAIL_ADDRESS")
	var email string = strings.ToLower(string(s.Personal.FirstName[0])) + "." + strings.ToLower(s.Personal.LastName) + addr
	if offset > 0 && offset < len([]rune(s.Personal.FirstName))-1 {
		email = lastEmail[:offset] + strings.ToLower(string(s.Personal.FirstName[offset])) + lastEmail[offset:] + addr
	}
	if offset == len([]rune(s.Personal.FirstName))-1 {
		email = strings.ToLower(s.Personal.FirstName) + "." + strings.ToLower(s.Personal.LastName) + addr
	}
	if offset > len([]rune(s.Personal.FirstName))-1 {
		email = strings.ToLower(s.Personal.FirstName) + "." + strings.ToLower(s.Personal.LastName) + strconv.Itoa(offset-len([]rune(s.Personal.FirstName))) + addr
	}
	return email
}

func (s *Student) ComparePasswords(password string) bool { //True: passwords match, False: no match
	valid := bcrypt.CompareHashAndPassword([]byte(s.Account.Password), []byte(password))
	return valid == nil
}

func (s *Student) CheckPasswordStrength(password string) bool {

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

func (s *Student) GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
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
