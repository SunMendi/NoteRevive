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
     token , err:=h.authService.LoginUser(req.Name,req.Email,req.TimeZone)
	 if err != nil {
		 ctx.JSON(http.StatusInternalServerError, gin.H{
			 "message":"err",
		 })
		 return
	 }

	 ctx.JSON(http.StatusOK, gin.H{
		 "message":token,
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

// GetProfile demonstrates how to access authenticated user information
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user information set by the middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userEmail, _ := c.Get("user_email")

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"user_email": userEmail,
		"message":    "Profile accessed successfully",
	})
}

// UpdateProfile demonstrates how to use authenticated user for updates
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the userID for database operations
	c.JSON(http.StatusOK, gin.H{
		"message":  "Profile updated successfully",
		"user_id": userID.(uint), // Type assertion since we know it's uint
	})
}