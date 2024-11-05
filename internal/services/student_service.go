package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
    s:= &StudentService{
        students: make(map[int]models.Student),
        nextID:   1,
    }
    // Loading initial data from JSON file
    s.LoadInitialData("data/students.json") 
    return s
}
// Loading initial data of students from a JSON file into the in-memory store
func (s *StudentService) LoadInitialData(filePath string) {
    fileData, err := os.ReadFile(filePath)
    if err != nil {
        log.Fatalf("Failed to read data file: %v", err)
    }

    var students []models.Student
    if err := json.Unmarshal(fileData, &students); err != nil {
        log.Fatalf("Failed to parse data file: %v", err)
    }

    s.mu.Lock()
    defer s.mu.Unlock()
    for _, student := range students {
        s.students[student.ID] = student
        if student.ID >= s.nextID {
            s.nextID = student.ID + 1
        }
    }
}

// writing the current students map to the JSON file
func (s *StudentService) saveToFile(filePath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Convert the map to a slice for JSON serialization
    students := make([]models.Student, 0, len(s.students))
    for _, student := range s.students {
        students = append(students, student)
    }

    data, err := json.MarshalIndent(students, "", "  ")
    if err != nil {
        return err
    }

    if err := os.WriteFile(filePath, data, 0644); err != nil {
        return err
    }

    return nil
}

func (s *StudentService) CreateStudent(student models.Student) (models.Student, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    student.ID = s.nextID
    s.students[student.ID] = student
    s.nextID++

    // Saving the updated student list to the file
    if err := s.saveToFile("data/students.json"); err != nil {
        log.Printf("Failed to save data: %v", err)
        return models.Student{}, err
    }

    return student, nil
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
