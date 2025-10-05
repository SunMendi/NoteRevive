package auth

import (
	"time"

	"gorm.io/gorm"
	"notemind/internal/note"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`

	Birthday time.Time `json:"birthday"`
	Gender   string    `json:"gender"`
	Timezone string    `json:"timezone"`

	Notes []note.Note `json:"notes" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
