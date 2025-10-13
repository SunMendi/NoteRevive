package auth

import "github.com/gin-gonic/gin"

func SetUpRoutes(router *gin.Engine, authHandler *AuthHandler) {
	v1 := router.Group("/api/v1")
	{

		v1.POST("/auth/user",authHandler.CreateUser)
		// Protected route - requires authentication
		v1.POST("/auth/send-daily-summary", authHandler.SendDailySummary)
	}

	// Protected routes that require authentication
	protected := router.Group("/api/v1/protected")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/profile", authHandler.GetProfile)
		protected.PUT("/profile", authHandler.UpdateProfile)
	}
}
