package routes

import (
	"school-management/controllers"

	"github.com/gofiber/fiber/v2"
)

type Person struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

func Setup(app *fiber.App) {
	// API Handling
	app.Post("/api/student/enroll", controllers.Enroll)
	app.Post("/api/student/updateName", controllers.UpdateStudentName)
	app.Post("/api/student/updateGradeLevel", controllers.UpdateStudentGradeLevel)
	app.Post("/api/student/updateHomeroom", controllers.UpdateStudentHomeroom)
	app.Post("/api/student/updateLocker", controllers.UpdateStudentLocker)
	app.Post("/api/studnet/updateYOG", controllers.UpdateStudentYOG)
	app.Post("/api/studnet/updateContacts", controllers.UpdateStudentContacts)
	app.Post("/api/student/updatePassword", controllers.UpdateStudentPassword)
	app.Post("/api/student/updateAddress", controllers.UpdateStudentAddress)
	app.Post("/api/student/updatePhoto", controllers.UpdateStudentPhoto)
	app.Post("/api/student/updateEmail", controllers.UpdateStudentEmail)

	app.Post("/api/teacher/register", controllers.RegisterTeacher)
	app.Post("/api/teacher/updatePassword", controllers.UpdateTeacherPassword)
	app.Post("/api/teacher/updateAddress", controllers.UpdateTeacherAddress)
	app.Post("/api/teacher/updatePhoto", controllers.UpdateTeacherPhoto)
	app.Post("/api/teacher/updateName", controllers.UpdateTeacherName)
	app.Post("/api/teacher/updateHomeroom", controllers.UpdateTeacherHomeroom)
	app.Post("/api/teacher/updateEmail", controllers.UpdateTeacherEmail)
}
