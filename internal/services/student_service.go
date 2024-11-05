package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"student-api/config"
	"student-api/internal/models"
	"sync"

	"github.com/go-resty/resty/v2"
)

type StudentService struct {
    students map[int]models.Student
    mu       sync.Mutex
    nextID   int
}

func NewStudentService() *StudentService {
    return &StudentService{
        students: make(map[int]models.Student),
        nextID:   1,
    }
}

func (s *StudentService) CreateStudent(student models.Student) models.Student {
    s.mu.Lock()
    defer s.mu.Unlock()
    student.ID = s.nextID
    s.students[student.ID] = student
    s.nextID++
    return student
}

func (s *StudentService) GetAllStudents() []models.Student {
    s.mu.Lock()
    defer s.mu.Unlock()
    studentList := make([]models.Student, 0, len(s.students))
    for _, student := range s.students {
        studentList = append(studentList, student)
    }
    return studentList
}

func (s *StudentService) GetStudentByID(id int) (models.Student, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    student, exists := s.students[id]
    if !exists {
        return models.Student{}, errors.New("student not found")
    }
    return student, nil
}

func (s *StudentService) UpdateStudent(id int, updatedStudent models.Student) (models.Student, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, exists := s.students[id]
    if !exists {
        return models.Student{}, errors.New("student not found")
    }
    updatedStudent.ID = id
    s.students[id] = updatedStudent
    return updatedStudent, nil
}

func (s *StudentService) DeleteStudent(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, exists := s.students[id]
    if !exists {
        return errors.New("student not found")
    }
    delete(s.students, id)
    return nil
}

func (s *StudentService) GenerateStudentSummary(id int) (string, error) {
    student, err := s.GetStudentByID(id)
    if err != nil {
        return "", err
    }

    client := resty.New()
    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetBody(map[string]interface{}{
            "prompt": fmt.Sprintf("Provide a summary for a student named %s, aged %d, email: %s.", student.Name, student.Age, student.Email),
        }).
        Post(config.OllamaAPIURL)

    if err != nil || resp.StatusCode() != http.StatusOK {
        return "", errors.New("failed to generate summary")
    }

    var result map[string]interface{}
    json.Unmarshal(resp.Body(), &result)

    summary, ok := result["summary"].(string)
    if !ok {
        return "", errors.New("invalid response from Ollama")
    }
    return summary, nil
}
