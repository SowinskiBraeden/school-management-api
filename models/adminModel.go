package models

import (
	"math/rand"
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
	Password     string             `json:"-" validate:"min=10,max=32"`
	TempPassword bool               `json:"temppassword"`
	AID          string             `json:"aid"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
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

	var hasUpper bool = true
	for _, r := range password {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			hasUpper = false
		}
	}

	var hasLower bool = true
	for _, r := range password {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			hasLower = false
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