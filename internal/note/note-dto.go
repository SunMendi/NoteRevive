package note

type CreateNoteDTO struct {
	Title   string `form:"title"`
	Content string `form:"content"`
}

// Add this to your existing DTO file:

type UpdateNoteDTO struct {
	Title   string `form:"title"`
	Content string `form:"content"`
}
