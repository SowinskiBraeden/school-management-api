package models

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*?,.`~"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

type Teacher struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"firstname" validate:"required"`
	MiddleName   string             `json:"middlename"`
	LastName     string             `json:"lastname" validate:"required"`
	Email        string             `json:"email" validate:"required"`
	SchoolEmail  string             `json:"schoolemail"`
	Password     string             `json:"-" validate:"min=10,max=32"`
	TempPassword bool               `json:"temppassword"`
	TID          string             `json:"tid"` // Teacher ID
	Homeroom     string             `json:"homeroom"`
	Address      string             `json:"address"`
	City         string             `json:"city"`
	Province     string             `json:"province"`
	Postal       string             `json:"postal"`
	DOB          string             `json:"dob" validate:"required"`
	Photo        string             `json:"string"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}

func (t *Teacher) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (t *Teacher) GenerateSchoolEmail() string {
	var email string = strings.ToLower(string(t.FirstName[0])) + "_" + strings.ToLower(t.LastName) + "@surreyschools.ca"
	// Add check to see if email already exists
	return email
}

func (t *Teacher) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
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

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper {
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
