package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"student-api/internal/models"
	"student-api/internal/services"
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
    var student models.Student
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    createdStudent := h.service.CreateStudent(student)
    utils.RespondWithJSON(w, http.StatusCreated, createdStudent)
}

// Handler to get all students
func (h *StudentHandler) GetAllStudents(w http.ResponseWriter, r *http.Request) {
    students := h.service.GetAllStudents()
    utils.RespondWithJSON(w, http.StatusOK, students)
}

// Handler to get a student by ID
func (h *StudentHandler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
        return
    }
    student, found := h.service.GetStudentByID(id)
    if found  != nil {
        utils.RespondWithError(w, http.StatusNotFound, "Student not found")
        return
    }
    utils.RespondWithJSON(w, http.StatusOK, student)
}

// Handler to update a student by ID
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
        return
    }

    var updatedData models.Student
    if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    updatedStudent, found := h.service.UpdateStudent(id, updatedData)
    if found != nil {
        utils.RespondWithError(w, http.StatusNotFound, "Student not found")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, updatedStudent)
}

// Handler to delete a student by ID
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
        return
    }

    deleted := h.service.DeleteStudent(id)
    if deleted != nil {
        utils.RespondWithError(w, http.StatusNotFound, "Student not found")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}

// Handler to get a summary of a student by ID
func (h *StudentHandler) GetStudentSummary(w http.ResponseWriter, r *http.Request) {
    idParam := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid student ID")
        return
    }

    summary, err := h.service.GenerateStudentSummary(id)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error generating summary")
        return
    }

    utils.RespondWithJSON(w, http.StatusOK, summary)
}
