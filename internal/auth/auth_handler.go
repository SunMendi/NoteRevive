package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
type AuthHandler struct {
	 authService AuthService 
}
func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}


func(h *AuthHandler) CreateUser(ctx *gin.Context) {
	 var req UserCreateRequest
	 if err := ctx.ShouldBindJSON(&req) ; err != nil {
		 ctx.JSON(http.StatusBadRequest, gin.H{
			 "message":"data is missing for this request",
		 })
		 return 
	 }
     err:=h.authService.CreateUser(req.Name,req.Email,req.TimeZone)
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{
			 "message":"err",
		 })
		 return
	 }

	 ctx.JSON(http.StatusOK, gin.H{
		 "message":"successfuly created user",
	 })
}




func (h *AuthHandler) SendDailySummary(c *gin.Context) {
	err := h.authService.SendDailySummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send daily summary: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daily summary sent successfully",
	})
}