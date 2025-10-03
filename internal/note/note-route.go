package note

import "github.com/gin-gonic/gin"

func SetUpRoutes(router *gin.Engine, notehandler *NoteHandler) {
	v1 := router.Group("/api/v1")
	v1.POST("/notes", notehandler.CreateNote)
	v1.PUT("/notes/:id", notehandler.UpdateNote)
	v1.GET("/notes/:id", notehandler.GetOneNote)
}
