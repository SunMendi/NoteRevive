package auth

import "github.com/gin-gonic/gin"

func SetUpRoutes(router *gin.Engine, authHandler *AuthHandler) {
	v1 := router.Group("/api/v1")
	{

		v1.POST("/auth/user",authHandler.CreateUser)
		v1.POST("/auth/send-daily-summary", authHandler.SendDailySummary)
	}
}
