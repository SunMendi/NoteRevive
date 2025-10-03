package main

import (
	"log"
	"notemind/database"
	"notemind/internal/auth"
	"notemind/internal/llm"
	"notemind/internal/note"
	"notemind/internal/voice"

	"github.com/gin-gonic/gin"
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

	note.SetUpRoutes(router, notehandler)
	auth.SetUpRoutes(router, authHandler)

	router.Run(":8080")

}
