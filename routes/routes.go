package routes

import (
	"school-management/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// API Handling

	// Student Authentication Handler
	app.Get("/api/student", controllers.Student)
	app.Post("/api/student/enroll", controllers.Enroll)
	app.Post("/api/student/login", controllers.StudentLogin)

	// Update Student Handler
	app.Post("/api/student/updateName", controllers.UpdateStudentName)
	app.Post("/api/student/updateGradeLevel", controllers.UpdateStudentGradeLevel)
	app.Post("/api/student/updateHomeroom", controllers.UpdateStudentHomeroom)
	app.Post("/api/student/updateLocker", controllers.UpdateStudentLocker)
	app.Post("/api/studnet/updateYOG", controllers.UpdateStudentYOG)
	app.Post("/api/studnet/addContact", controllers.AddStudentContact)
	app.Post("/api/student/removeContact", controllers.RemoveStudentContact)
	app.Post("/api/student/updatePassword", controllers.UpdateStudentPassword)
	app.Post("/api/student/resetPassword", controllers.ResetStudentPassword)
	app.Post("/api/student/updateAddress", controllers.UpdateStudentAddress)
	app.Post("/api/student/updatePhoto", controllers.UpdateStudentPhoto)
	app.Post("/api/student/updateEmail", controllers.UpdateStudentEmail)

	// Student Contact Handler
	app.Post("/api/student/createContact", controllers.CreateContact)
	app.Post("/api/student/contact/updateName", controllers.UpdateContactName)
	app.Post("/api/student/contact/updateAddress", controllers.UpdateContactAddress)
	app.Post("/api/student/contact/updateHomePhone", controllers.UpdateContactHomePhone)
	app.Post("/api/student/contact/updateWorkPhone", controllers.UpdateContactWorkPhone)
	app.Post("/api/student/contact/updateEmail", controllers.UpdateContactEmail)
	app.Post("/api/student/contact/updatePriority", controllers.UpdateContactPriority)
	app.Post("/api/student/contact/deleteContact", controllers.DeleteContact)

	// Teacher Login Handling
	app.Get("/api/teacher", controllers.Teacher)
	app.Post("/api/teacher/register", controllers.RegisterTeacher)
	app.Post("/api/teacher/login", controllers.TeacherLogin)

	// Update Teacher Handler
	app.Post("/api/teacher/updatePassword", controllers.UpdateTeacherPassword)
	app.Post("/api/teacher/updateAddress", controllers.UpdateTeacherAddress)
	app.Post("/api/teacher/updatePhoto", controllers.UpdateTeacherPhoto)
	app.Post("/api/teacher/updateName", controllers.UpdateTeacherName)
	app.Post("/api/teacher/updateHomeroom", controllers.UpdateTeacherHomeroom)
	app.Post("/api/teacher/updateEmail", controllers.UpdateTeacherEmail)
	app.Post("/api/teacher/resetPassword", controllers.ResetTeacherPassword)

	// General Routes
	app.Post("/api/logout", controllers.Logout)

	// Admin Login Handling
	app.Get("/api/admin", controllers.Admin)
	app.Post("/api/admin/create", controllers.CreateAdmin)
	app.Post("/api/admin/login", controllers.AdminLogin)

	// General Command Handling
	app.Post("/api/admin/updateLockerCombo", controllers.UpdateLockerCombo)
	app.Post("/api/admin/renableStudent", controllers.RemoveStudentsDisabled)
}
