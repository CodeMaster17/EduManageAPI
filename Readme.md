## Folder structure

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
