package routes

import (
	"github.com/SowinskiBraeden/school-management-api/controllers"
	"github.com/SowinskiBraeden/school-management-api/controllers/update"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Detect if system is new and needs default admin
	controllers.NewSystem()

	// API Handling
	var routerPrefix string = "/api/v1"

	// API check
	app.Get(routerPrefix+"/status", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "the API is active",
		})
	})

	// Student Authentication Handler
	app.Get(routerPrefix+"/student", controllers.Student)
	app.Post(routerPrefix+"/student/enroll", controllers.Enroll)
	app.Post(routerPrefix+"/student/login", controllers.StudentLogin)

	// Update Student Handler
	app.Post(routerPrefix+"/student/updateName", update.UpdateStudentName)
	app.Post(routerPrefix+"/student/updateGradeLevel", update.UpdateStudentGradeLevel)
	app.Post(routerPrefix+"/student/updateHomeroom", update.UpdateStudentHomeroom)
	app.Post(routerPrefix+"/student/updateLocker", update.UpdateStudentLocker)
	app.Post(routerPrefix+"/studnet/updateYOG", update.UpdateStudentYOG)
	app.Post(routerPrefix+"/studnet/addContact", update.AddStudentContact)
	app.Post(routerPrefix+"/student/removeContact", update.RemoveStudentContact)
	app.Post(routerPrefix+"/student/updatePassword", update.UpdateStudentPassword)
	app.Post(routerPrefix+"/student/resetPassword", update.ResetStudentPassword)
	app.Post(routerPrefix+"/student/updateAddress", update.UpdateStudentAddress)
	app.Post(routerPrefix+"/student/updatePhoto", update.UpdateStudentPhoto)
	app.Post(routerPrefix+"/student/updateEmail", update.UpdateStudentEmail)

	// Student Contact Handler
	app.Post(routerPrefix+"/contact/createContact", controllers.CreateContact)
	app.Post(routerPrefix+"/contact/updateName", update.UpdateContactName)
	app.Post(routerPrefix+"/contact/updateAddress", update.UpdateContactAddress)
	app.Post(routerPrefix+"/contact/updateHomePhone", update.UpdateContactHomePhone)
	app.Post(routerPrefix+"/contact/updateWorkPhone", update.UpdateContactWorkPhone)
	app.Post(routerPrefix+"/contact/updateEmail", update.UpdateContactEmail)
	app.Post(routerPrefix+"/contact/updatePriority", update.UpdateContactPriority)
	app.Post(routerPrefix+"/contact/deleteContact", controllers.DeleteContact)

	// Teacher Authentication Handler
	app.Get(routerPrefix+"/teacher", controllers.Teacher)
	app.Post(routerPrefix+"/teacher/register", controllers.RegisterTeacher)
	app.Post(routerPrefix+"/teacher/login", controllers.TeacherLogin)

	// Teacher Update Handler
	app.Post(routerPrefix+"/teacher/updatePassword", update.UpdateTeacherPassword)
	app.Post(routerPrefix+"/teacher/updateAddress", update.UpdateTeacherAddress)
	app.Post(routerPrefix+"/teacher/updatePhoto", update.UpdateTeacherPhoto)
	app.Post(routerPrefix+"/teacher/updateName", update.UpdateTeacherName)
	app.Post(routerPrefix+"/teacher/updateHomeroom", update.UpdateTeacherHomeroom)
	app.Post(routerPrefix+"/teacher/updateEmail", update.UpdateTeacherEmail)
	app.Post(routerPrefix+"/teacher/resetPassword", update.ResetTeacherPassword)

	// General Routes
	app.Post(routerPrefix+"/logout", controllers.Logout)

	// Admin Login Handling
	app.Get(routerPrefix+"/admin", controllers.Admin)
	app.Post(routerPrefix+"/admin/create", controllers.CreateAdmin)
	app.Post(routerPrefix+"/admin/login", controllers.AdminLogin)

	// Admin Update Handler
	app.Post(routerPrefix+"/admin/updateName", update.UpdateAdminName)
	app.Post(routerPrefix+"/admin/updateEmail", update.UpdateAdminEmail)
	app.Post(routerPrefix+"/admin/updatePassword", update.UpdateAdminPassword)

	// General Command Handling
	app.Post(routerPrefix+"/admin/updateLockerCombo", update.UpdateLockerCombo)
	app.Post(routerPrefix+"/admin/enableStudent", update.RemoveStudentsDisabled)
	app.Post(routerPrefix+"/admin/enableTeacher", update.RemoveTeachersDisabled)

	// Delete Handler
	app.Post(routerPrefix+"/remove/student", controllers.RemoveStudent)
	app.Post(routerPrefix+"/remove/teacher", controllers.RemoveTeacher)
	app.Post(routerPrefix+"/remove/admin", controllers.RemoveAdmin)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}
