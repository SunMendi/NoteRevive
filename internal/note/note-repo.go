package note

import (
	"gorm.io/gorm"
)

type NoteRepo interface {
	Create(note *Note) error
	CreateImg(noteImage *NoteImage) error
	Update(note *Note) error
	Delete(id uint) error 
	GetByID(id uint) (*Note, error)
	DeleteImagesByNoteID(id uint) error
}

type noterepo struct {
	db *gorm.DB
}

func NewNoteRepo(db *gorm.DB) NoteRepo {
	return &noterepo{db: db}
}

func (r *noterepo) Create(note *Note) error {
	return r.db.Create(note).Error
}

func (r *noterepo) CreateImg(noteImage *NoteImage) error {
	return r.db.Create(noteImage).Error
}

func (r *noterepo) GetByID(id uint) (*Note, error) {
	var note Note
	if err := r.db.Preload("Images").First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noterepo) Update(note *Note) error {
	return r.db.Save(note).Error
}

func (r *noterepo) DeleteImagesByNoteID(noteID uint) error {
	return r.db.Where("note_id = ?", noteID).Delete(&NoteImage{}).Error
}

func(r *noterepo) Delete(id uint) error {
	return r.db.Delete(&Note{}, id).Error 
}