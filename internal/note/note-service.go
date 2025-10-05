package note

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"
	"log"

	"notemind/internal/llm"
	"notemind/internal/voice"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type NoteService interface {
	CreateNote(userID uint, title string, content string, imageFile *multipart.FileHeader) (*Note, error)
	CreateVoiceNote(userID uint, audioFile *multipart.FileHeader, title string, imageFile *multipart.FileHeader) (*Note, error)
	UpdateNote(noteID uint, userID uint, title, content string, imageFile *multipart.FileHeader) error
	GetOneNote(id uint) (*Note, error)
}

type noteService struct {
	repo        NoteRepo
	llmservice  *llm.LLMService
	transcriber voice.Transcriber
}

func NewNoteService(repo NoteRepo, llmService *llm.LLMService, transcriber voice.Transcriber) NoteService {
	return &noteService{
		repo:        repo,
		llmservice:  llmService,
		transcriber: transcriber,
	}
}

// ...existing code...

func (s *noteService) handleImageUpload(noteID uint, imageFile *multipart.FileHeader) error {
	src, err := imageFile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return errors.New("cloudinary config failed")
	}

	result, err := cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			ResourceType: "image",
			Folder:       "notes",
		},
	)
	if err != nil {
		return errors.New("upload failed: " + err.Error())
	}

	noteImage := &NoteImage{
		NoteID:     noteID,
		ImageURL:   result.SecureURL,
		UploadedAt: time.Now(),
	}

	return s.repo.CreateImg(noteImage)
}

func (s *noteService) CreateNote(userID uint, title, content string, imageFile *multipart.FileHeader) (*Note,error) {
	if userID == 0 {
		return nil, errors.New("user ID can not become zero")
	}
	noteText := fmt.Sprintf("Title: %s\nContent: %s", title, content)
	log.Println(noteText)
	summary, err := s.llmservice.GenerateNoteSummary(noteText)

	if err != nil {
		// If summary generation fails, continue without summary
		summary = "Summary generation failed"
	}
	note := &Note{
		UserID:    userID,
		Title:     title,
		Content:   content,
		Summary:   summary,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(note); err != nil {
		return nil,err
	}

	if imageFile != nil {
		if err := s.handleImageUpload(note.ID, imageFile); err != nil {
			return nil, err
		}
	}
	return note, nil
}

func (s *noteService) CreateVoiceNote(userID uint, audioFile *multipart.FileHeader, title string, imageFile *multipart.FileHeader) (*Note, error) {
	if userID == 0 {
		return nil,errors.New("user ID cannot be zero")
	}
	if audioFile == nil {
		return nil,errors.New("audio file is required")
	}

	audio, err := audioFile.Open()
	if err != nil {
		return nil,err
	}

	transcript, err := s.transcriber.Transcribe(context.Background(), audio)
	if err != nil {
		return nil,err
	}

	noteTitle := title
	if noteTitle == "" {
		noteTitle = audioFile.Filename
	}

	return s.CreateNote(userID, noteTitle, transcript, imageFile)
}

// Add this to your existing NoteService interface:

func (s *noteService) UpdateNote(noteID uint, userID uint, title string, content string, imageFile *multipart.FileHeader) error {
	if userID == 0 {
		return errors.New("user ID cannot be zero")
	}

	if noteID == 0 {
		return errors.New("note ID cannot be zero")
	}

	// STEP 1: Get existing note and verify ownership
	existingNote, err := s.repo.GetByID(noteID)
	if err != nil {
		return errors.New("note not found")
	}

	// Check if user owns this note
	if existingNote.UserID != userID {
		return errors.New("unauthorized: you don't own this note")
	}

	// STEP 2: Generate new summary if content changed
	var summary string
	if content != existingNote.Content || title != existingNote.Title {
		noteText := fmt.Sprintf("Title: %s\nContent: %s", title, content)
		summary, err = s.llmservice.GenerateNoteSummary(noteText)
		if err != nil {
			// If summary generation fails, keep old summary
			summary = existingNote.Summary
		}
	} else {
		// Content didn't change, keep existing summary
		summary = existingNote.Summary
	}

	// STEP 3: Update note fields
	if title != "" {
		existingNote.Title = title
	}
	if content != "" {
		existingNote.Content = content
	}
	if summary != "" {
		existingNote.Summary = summary
	}
	existingNote.UpdatedAt = time.Now()

	// STEP 4: Save updated note
	if err := s.repo.Update(existingNote); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	// STEP 5: Handle image update if new image provided
	if imageFile != nil {
		// Delete old images first
		if err := s.repo.DeleteImagesByNoteID(noteID); err != nil {
			return fmt.Errorf("failed to delete old images: %w", err)
		}

		// Upload new image
		if err := s.handleImageUpload(noteID, imageFile); err != nil {
			return fmt.Errorf("failed to upload new image: %w", err)
		}
	}
	return nil
}

func (s *noteService) GetOneNote(id uint) (*Note, error) {
	if id == 0 {
		return nil, errors.New("zero value")
	}
	res, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil

}
