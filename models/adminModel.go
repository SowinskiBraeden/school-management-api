package models

import (
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (a *Admin) GenerateSchoolEmail() string {
	addr := os.Getenv("SYSTEM_EMAIL_ADDRESS")
	var email string = strings.ToLower(a.LastName) + "_" + strings.ToLower(string(a.FirstName[0])) + addr
	// Add check to see if email already exists
	return email
}

func (s *Admin) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (s *Admin) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (s *Admin) CheckPasswordStrength(password string) bool {

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

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper {
		return true
	} else {
		return false
	}
}

func (s *Admin) GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
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
