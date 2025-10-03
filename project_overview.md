NoteMind - Technical Overview
1. High-Level Summary
Project Purpose
NoteMind is a comprehensive note-taking application that combines traditional text-based note creation with modern AI-powered features. The application allows users to create, update, and manage notes while leveraging artificial intelligence for automatic summarization and voice-to-text transcription capabilities. The target users are individuals who want an intelligent note-taking system that can process both text and audio inputs while providing AI-generated summaries for better organization and quick reference.

Core Features
Text Note Management: Create, read, update, and delete text-based notes with titles and content
Voice Note Transcription: Convert audio files to text using Deepgram's speech-to-text API
AI-Powered Summarization: Automatically generate concise summaries of notes using Google's Generative AI (Gemini)

Image Attachment: Upload and attach images to notes using Cloudinary for cloud storage
User Authentication & Authorization: Secure user management with ownership validation for notes
2. Technology Stack
Languages
Go: Primary backend language for the entire application
Backend
Framework: Custom Go HTTP server (likely using standard library or minimal framework)
Database ORM/Driver: Custom repository pattern implementation
Key Libraries:
Google Generative AI Go SDK (github.com/google/generative-ai-go/genai)
Deepgram Go SDK (github.com/deepgram/deepgram-go-sdk/v3)
Cloudinary Go SDK (github.com/cloudinary/cloudinary-go/v2)

Database
Database System: SQL-based database (likely PostgreSQL) with migration support
Migrations: Database versioning with up/down migration files
External Services
AI/ML: Google Generative AI (Gemini 1.5 Flash) for text summarization
Speech-to-Text: Deepgram API for audio transcription
Image Storage: Cloudinary for image upload and management
Tools & DevOps
Dependency Management: Go modules (go.mod)
Environment Configuration: .env file for environment variables
Hot Reload: Temporary build directory (tmp) suggesting development tooling
3. Project Architecture & Structure
Architecture
The project follows a Clean Architecture/Layered Architecture pattern with clear separation of concerns:

Handler Layer: HTTP request handling and routing
Service Layer: Business logic implementation
Repository Layer: Data access abstraction
External Services: Third-party integrations (AI, transcription, cloud storage)
Directory Overview - 
notemind/
├── .env                          # Environment variables
├── main.go                       # Application entry point
├── go.mod                        # Go module dependencies
├── database/
│   └── database.go              # Database connection and configuration
├── internal/
│   ├── auth/                    # Authentication and authorization
│   │   ├── auth_dtos.go        # Data transfer objects
│   │   └── auth_handler.go     # HTTP handlers
│   ├── llm/                     # AI/LLM integration
│   │   └── llm-service.go      # Generative AI service
│   ├── note/                    # Core note functionality
│   │   └── note-service.go     # Note business logic
│   ├── summary/                 # Summary-related features
│   └── voice/                   # Voice transcription
│       └── deepgram-client.go  # Deepgram integration
├── migrations/                  # Database migration files
│   ├── 000001_*.up.sql         # User model migrations
│   ├── 000002_*.up.sql         # Note and Image model migrations
│   └── 000003_*.up.sql         # Summary field addition
└── tmp/     
                    # Temporary build files
Key Directory Explanations
internal: Contains all internal application packages following Go conventions
note: Core note management functionality including CRUD operations
auth: User authentication, authorization, and session management
llm: Integration with Google's Generative AI for text summarization
voice: Audio processing and transcription using Deepgram
migrations: Database schema evolution with versioned SQL migration files
database: Database connection setup and configuration
4. Core Logic and Data Flow
Entry Point
The application starts from main.go, which likely initializes the HTTP server, database connections, and external service clients.

Data Flow
Note Creation: HTTP request → auth_handler.go (authentication) → note-service.go → llm-service.go (summary generation) → Repository → Database
Voice Note Processing: Audio upload → deepgram-client.go (transcription) → note-service.go → llm-service.go → Repository → Database
Image Upload: Image file → note-service.go → Cloudinary API → Database (URL storage)

API Routes (Inferred)
POST /notes - Create a new text note
POST /notes/voice - Create a note from audio file
PUT /notes/:id - Update an existing note
GET /notes/:id - Retrieve a specific note
Authentication endpoints in auth handler
Authentication
User authentication is handled through the auth package with ownership validation ensuring users can only access their own notes. The UpdateNote function includes authorization checks to verify note ownership.

Key Files
note-service.go: Core business logic for note operations including CRUD, image handling, and AI integration
llm-service.go: AI summarization service using Google's Generative AI with custom prompting
deepgram-client.go: Voice transcription service integrating with Deepgram's speech-to-text API

5. Setup & How to Run
Prerequisites
Go: Version 1.19+ (based on go.mod requirements)
Database: PostgreSQL or compatible SQL database
