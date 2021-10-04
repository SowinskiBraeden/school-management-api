package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID         primitive.ObjectID `bson:"_id"`
	FirstName  string             `json:"firstname" validate:"required"`
	MiddleName string             `json:"middlename"`
	LastName   string             `json:"lastname" validate:"required"`
	Age        int                `json:"age" validate:"required"`
	GradeLevel int                `json:"gradelevel" validate:"required"`
	Email      string             `json:"email" validate:"required"`
	Password   string             `json:"password" validate:"required,min=8,max=32"`
	SID        string             `json:"sid"` // Student ID
	PED        string             `json:"ped"` // Personal Education Number
	Homeroom   string             `json:"homeroom"`
	Locker     string             `json:"locker"`
	YOG        string             `json:"yog"` // Year of Graduation
	Address    string             `json:"address"`
	City       string             `json:"city"`
	Province   string             `json:"province"`
	PC         string             `json:"pc"` // Postal Code
	DOB        string             `json:"dob" validate:"required"`
	Photo      string             `json:"photo"`
	Contacts   []Contact          `json:"contacts"`
}
