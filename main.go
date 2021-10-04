package main

import (
	"fmt"
	"school-management/models"
)

func main() {
	var student models.Student
	student.FirstName = "Braeden"
	student.LastName = "Sowinski"
	student.Age = 16
	student.DOB = "2005-04-26" //year-month-day
	fmt.Print(student)
}
