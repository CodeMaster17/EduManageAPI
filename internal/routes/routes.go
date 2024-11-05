package routes

import (
	"student-api/internal/handlers"
	"student-api/internal/services"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
    studentService := services.NewStudentService()
    studentHandler := handlers.NewStudentHandler(studentService)

    r.Post("/students", studentHandler.CreateStudent)
    r.Get("/students", studentHandler.GetAllStudents)
    r.Get("/students/{id}", studentHandler.GetStudentByID)
    r.Put("/students/{id}", studentHandler.UpdateStudent)
    r.Delete("/students/{id}", studentHandler.DeleteStudent)
    r.Get("/students/{id}/summary", studentHandler.GetStudentSummary)
}
