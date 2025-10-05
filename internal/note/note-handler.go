package note

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	noteService NoteService
}

func NewNoteHandler(noteService NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) CreateNote(ctx *gin.Context) {
	// Get user_id from middleware context with proper type assertion
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert interface{} to uint
	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid user ID type"})
		return
	}

	var req CreateNoteDTO
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ADD VALIDATION - Prevent empty notes

	imageFile, err := ctx.FormFile("image")
	if err != nil {
		log.Println("image file is missing")
	}
	
	audioFile, err := ctx.FormFile("audio")
	if err == nil {
		// Create voice note - declare note variable here
		note, err := h.noteService.CreateVoiceNote(userID, audioFile, req.Title, imageFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "Voice note created successfully",
			"note": note,
		})
		return
	}

	// Create regular note - declare note variable here
	note, err := h.noteService.CreateNote(userID, req.Title, req.Content, imageFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Note created successfully",
		"note": note,
	})
}


func (h *NoteHandler) UpdateNote(ctx *gin.Context) {

	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Convert interface{} to uint
	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid user ID type"})
		return
	}

	noteIDStr := ctx.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid note ID",
		})
		return
	}

	var req UpdateNoteDTO
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	imageFile, err := ctx.FormFile("image")
	if err != nil && err.Error() != "http: no such file" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Error reading image file",
		})
		return
	}
	// Replace with actual user from JWT middleware

	// STEP 5: Update note
	err = h.noteService.UpdateNote(uint(noteID), userID, req.Title, req.Content, imageFile)
	if err != nil {
		// Handle different error types
		if err.Error() == "note not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Note not found",
			})
			return
		}
		if err.Error() == "unauthorized: you don't own this note" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized access",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// STEP 6: Success response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Note updated successfully",
		"note_id": noteID,
	})
}

func (h *NoteHandler) GetOneNote(ctx *gin.Context) {
	id := ctx.Param("id")
	idStr, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid note ID",
		})
		return
	}
	if idStr == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id can not become zero",
		})
		return
	}
	res, err := h.noteService.GetOneNote(uint(idStr))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"note": res,
	})

}
