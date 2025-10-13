package note

import (
	"notemind/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine, notehandler *NoteHandler) {
	v1 := router.Group("/api/v1")
	v1.POST("/notes", middleware.AuthMiddleware(), notehandler.CreateNote)
	v1.PUT("/notes/:id", middleware.AuthMiddleware(),notehandler.UpdateNote)
	v1.GET("/notes/:id",middleware.AuthMiddleware(), notehandler.GetOneNote)
}
