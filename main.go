package main

import (
	"log"
	"notemind/database"
	"notemind/internal/auth"
	"notemind/internal/llm"
	"notemind/internal/note"
	"notemind/internal/voice"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {

	db, err := database.ConnectDB()

	if err != nil {
		log.Panic("Database initilization issue")
	}

	log.Printf("Database connected successfully")

	llmService, err := llm.NewLLMService()

	if err != nil {
		log.Println("LLM service initialization issue")
	}

	log.Printf("gemini integrate successfully")

	voiceClient, err := voice.NewdeepgramClient()

	if err != nil {
		log.Println("failed to init deepgram")
	}

	defer llmService.Close()

	gin.SetMode(gin.ReleaseMode)

	noteRepo := note.NewNoteRepo(db)
	authRepo := auth.NewAuthRepo(db, llmService)

	//log.Println(authRepo)

	noteService := note.NewNoteService(noteRepo, llmService, voiceClient)
	authService := auth.NewAuthService(authRepo) 




	notehandler := note.NewNoteHandler(noteService)
	authHandler := auth.NewAuthHandler(authService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://noterevive-production.up.railway.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	note.SetUpRoutes(router, notehandler)
	auth.SetUpRoutes(router, authHandler)

	router.Run(":8080")

}
