## EduManage API

 A simple go application made folllowing SOLID Principles and DRY Principle, that generates summary of students using Ollama LLM using LLama3.2 model.

## Folder structure

```
edumanageapi
├── cmd/                    # Main applications for the project
│   └── server/             # Main application entry point
│       └── main.go         # Main application file
── config/                   # Configuration files and utilities
│   └── config.go           
├── internal/               # Private application and library code
│   ├── handlers/           # HTTP handlers (Controllers)
│   │   └── student_handler.go      # Student-related HTTP handlers
│   ├── services/           # Business logic (Services)
│   │   └── student_service.go      # Student service layer
│   ├── models/             # Data models
│   │   └── student.go      # Student data model
│   ├── routes/         # Data access layer
│   │   └── routes.go # CRUD operations for students
│   ├── utils/              # Utility functions/helpers
│   │   └── response.go     # JSON response helper functions
└── README.md               # Project README file                # 
├── go.mod                  # Go module file
└── go.sum                  # Go dependencies checksum file

```


## Installation

1. Clone the repository
```
git clone https://github.com/CodeMaster17/EduManageAPI.git

```

2. Change the directory
```
cd EduManageAPI
```
3. Start the application
```
go run cmd/server/main.go
```
4. Install and run Ollama
```
1. Install Ollama on your localhost by following: https://github.com/ollama/ollama/blob/main/README.md#quickstart
2. Install Llama3 language model for Ollama
3. Make API Requests to localhost and pass the Student object to generate the summary. Do the prompt engineering to get the summary for the Student.

```