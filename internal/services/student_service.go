package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"student-api/config"
	"student-api/internal/types"

	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/rand"
)

// Message struct to use within request/response structures
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Student struct with ID, name, and other details
type Student struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Email   string  `json:"email"`
	Message Message `json:"message"`
}

// StudentService to manage students in-memory
type StudentService struct {
	students map[int]Student
}

var students []types.Student

// NewStudentService initializes the StudentService with an empty map and ID counter
func NewStudentService() *StudentService {
	return &StudentService{
		students: make(map[int]Student),
	}
}

// get all students
func GetAllStudents(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// --- CreateStudent creates a new student and adds it to the map ---
func (s *StudentService) CreateStudent(student types.Student) ([]types.Student, error) {

	student.ID = strconv.Itoa(rand.Intn(1000000))
	students = append(students, student)
	fmt.Println("Created student Inside: ", student)
	return students, nil
}

// get all the students
func (s *StudentService) GetAllStudents() ([]types.Student, error) {
	return students, nil
}

// GetStudentByID retrieves a student by ID
func (s *StudentService) GetStudentByID(id string) (types.Student, error) {
	for _, item := range students {
		if item.ID == id {
			return item, nil
		}
	}
	return types.Student{}, errors.New("student not found")
}

func (s *StudentService) UpdateStudent(id string, updatedData types.Student) (types.Student, error) {
	for index, item := range students {
		if item.ID == id {

			updatedData.ID = id
			students[index] = updatedData
			return updatedData, nil
		}
	}
	return types.Student{}, errors.New("student not found")
}

// DeleteStudent deletes a student by ID
func (s *StudentService) DeleteStudent(id string) error {
	for index, item := range students {
		if item.ID == id {
			students = append(students[:index], students[index+1:]...)
			return nil
		}
	}
	return errors.New("student not found")
}

// GenerateStudentSummary generates a summary for a student by ID
func (s *StudentService) GenerateStudentSummary(id string) (string, error) {
	student, err := s.GetStudentByID(id)
	if err != nil {
		return "", err
	}

	// Create the prompt for the Ollama API
	prompt := fmt.Sprintf("Generate fictional summary and characterstics of student with name %s and email %s", student.Name, student.Email)

	// Initialize Resty client
	client := resty.New()

	if config.OllamaAPIURL == "" {
		return "", errors.New("ollama API URL is not set")
	}
	fmt.Println("Program here 2")

	// Make the request to the Ollama API
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":  "llama3.2",
			"prompt": prompt,
			"format": "json",
			"stream": false,
		}).
		Post(config.OllamaAPIURL)

	// Check for errors and the response status code
	if err != nil {
		return "", errors.New("failed to generate summary: " + err.Error())
	}
	fmt.Println("Program here 2", resp)
	fmt.Println(resp.StatusCode())
	if resp.StatusCode() != http.StatusOK {
		fmt.Println("Program inside status code")
		return "", fmt.Errorf("ollama API returned status code %d", resp.StatusCode())
	}

	// Parse the response from the Ollama API
	var result map[string]interface{}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", errors.New("error parsing Ollama API response: " + err.Error())
	}

	// Extract the summary from the response
	responseString, ok := result["response"].(string)
	if !ok {
		return "", errors.New("invalid response format from Ollama API")
	}

	// Unmarshal the nested JSON string to extract the content
	var responseContent map[string]interface{}
	if err := json.Unmarshal([]byte(responseString), &responseContent); err != nil {
		return "", errors.New("error parsing nested response JSON: " + err.Error())
	}

	// Extract the actual content you want to return
	content, ok := responseContent["content"].(string)
	if !ok {
		return "", errors.New("invalid content format in nested response")
	}

	return content, nil
}
