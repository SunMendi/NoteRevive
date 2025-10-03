package note

import "time"

type Note struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	UserID  uint   `json:"user_id"` //foreign key
	Title   string `json:"title"`
	Content string `json:"content"`
	Summary string `json:"summary"`

	Images    []NoteImage `json:"images,omitempty" gorm:"foreignKey:NoteID"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type NoteImage struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	NoteID   uint   `json:"note_id"` // foreign key
	ImageURL string `json:"image_url"`

	UploadedAt time.Time `json:"uploaded_at"`
}
