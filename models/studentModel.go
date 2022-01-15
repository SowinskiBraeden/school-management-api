package models

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID           primitive.ObjectID `bson:"_id"`
	PersonalData struct {
		FirstName  string   `json:"firstname" validate:"required"`
		MiddleName string   `json:"middlename"`
		LastName   string   `json:"lastname" validate:"required"`
		Age        int      `json:"age" validate:"required"`
		Email      string   `json:"email" validate:"required"`
		Address    string   `json:"address"`
		City       string   `json:"city"`
		Province   string   `json:"province"`
		Postal     string   `json:"postal"`
		DOB        string   `json:"dob" validate:"required"`
		Photo      string   `json:"photo"`
		Contacts   []string `json:"contacts"` // List of contact ID's rather than contact object
	}
	SchoolData struct {
		GradeLevel  int    `json:"gradelevel" validate:"required"`
		SchoolEmail string `json:"schoolemail"`
		Password    string `json:"-" validate:"min=10,max=32"`
		SID         string `json:"sid"` // Student ID
		PEN         string `json:"ped"` // Personal Education Number
		Homeroom    string `json:"homeroom"`
		Locker      Locker `json:"locker"`
		YOG         int    `json:"yog"` // Year of Graduation
		Photo       string `json:"photo"`
	}
	AccountData struct {
		AccountDisabled bool `bson:"accountdisabled"`
		TempPassword    bool `json:"temppassword"`
		Attempts        int  `json:"attempts"` // login attempts max 5
	}
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (s *Student) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (s *Student) GenerateSchoolEmail() string {
	var email string = strings.ToLower(string(s.PersonalData.FirstName[0])) + "." + strings.ToLower(s.PersonalData.LastName) + "@surreyschools.ca"
	// Add check to see if email already exists
	return email
}

func (s *Student) ComparePasswords(password string) bool {
	valid := bcrypt.CompareHashAndPassword([]byte(s.SchoolData.Password), []byte(password))
	if valid != nil {
		return false
	}
	return true
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

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper {
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
