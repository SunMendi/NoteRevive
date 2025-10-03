package auth


type UserCreateRequest struct {

	 Name string `json:"name" binding:"required"`
	 Email string `json:"email" binding:"required,email"`
	 TimeZone string `json:"timezone" binding:"required"`

}