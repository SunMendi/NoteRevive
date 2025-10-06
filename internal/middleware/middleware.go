package middleware


import (
	"os"
	"errors"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}


// ValidateToken validates a JWT token string
func ValidateToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("SECRET_KEY not found")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// AuthMiddleware creates a Gin middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	 return func(c *gin.Context) {
		 authHeader:=c.GetHeader("Authorization")
		 if authHeader == "" {
			 c.JSON(401, gin.H{
				"errors":"Authorization Header Required",
			 })
			 log.Println("token is missing")
			 c.Abort()
			 return 
		 }
		 if !strings.HasPrefix(authHeader, "Bearer " ){
			 c.JSON(401, gin.H{"error": "Authorization header must start with Bearer"})
             c.Abort()
			 return 
		 }
		 tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		 claims, err := ValidateToken(tokenString)
         if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
            c.Abort()
            return
        }
		c.Set("user_id",claims.UserID)

		c.Next() 
	 }
}