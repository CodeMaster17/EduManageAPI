package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"student-api/internal/services"
	"student-api/internal/types"
	"student-api/internal/utils"

	"github.com/go-chi/chi/v5"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// Handler to create a new student
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student types.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdStudent, err := h.service.CreateStudent(student)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not create student")
		return
	}
	fmt.Println("Created student: ", createdStudent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdStudent)
}

// Handler to get all students
func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, _ := h.service.GetAllStudents()
	utils.RespondWithJSON(w, http.StatusOK, students)
}

// Handler to get a student by ID
func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	student, err := h.service.GetStudentByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Student not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, student)
}

// Handler to update a student by ID
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Get the ID from URL

	var updatedData types.Student
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call the service to update the student
	updatedStudent, err := h.service.UpdateStudent(id, updatedData)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Student not found")
		return
	}

	// Respond with the updated student data
	utils.RespondWithJSON(w, http.StatusOK, updatedStudent)
}

// Handler to delete a student by ID
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // Get the ID from URL

	// Call the service to delete the student
	err := h.service.DeleteStudent(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Student not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}

// Handler to get a summary of a student by ID
func (h *StudentHandler) GetStudentSummary(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	summary, err := h.service.GenerateStudentSummary(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating summary")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, summary)
}
