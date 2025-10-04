package auth

import (
	"errors"
	"time"
	"os"
	"log"

	"github.com/golang-jwt/jwt/v5"
)


type AuthService interface {
     LoginUser(name, email , timezone string) (string , error ) 
	 GenerateToken(userID uint, email string) (string , error)

	 SendDailySummary() error 
}

type authService struct {
	 repo AuthRepo 
}

type JWTClaims struct {
	 UserID uint `json:"user_id"`
	 Email string `json:"email"`
	 jwt.RegisteredClaims
}


func NewAuthService (repo AuthRepo) AuthService {
	 return &authService{repo: repo}
}

func(s *authService) LoginUser(name, email , timezone string) (string , error) {
	 if name=="" {
		 return "", errors.New("name is required")
	 }
	 if email=="" {
		 return "",errors.New("email is required")
	 }
	 if timezone=="" {
		 return "",errors.New("timezone is required")
	 }

	 _, err := s.repo.GetByEmail(email)
	 if err != nil {  
	 user:=&User{
		Name:name,
		Email: email,
		Timezone: timezone,
		CreatedAt: time.Now().UTC(),
	 }
	 err:=s.repo.Create(user)
	 if err != nil {
		log.Println("user creation mistake")
		 return "",err  
	 }
	 }
	 user, err := s.repo.GetByEmail(email)
	 if err != nil {
		log.Println("user does not exists")
		 return "", err 
	 }
	 res , err :=s.GenerateToken(user.ID, user.Email)
	 if err != nil {
		log.Println("token creation issue")
		 return "", err 
	 }
	 return res , nil 

}

func(s *authService) GenerateToken(userID uint , email string) (string , error) {
	 secretKey:= os.Getenv("SECRET_KEY")
	 if secretKey == "" {
		 return "",errors.New("secret key is required")
	 }
	 claims:= &JWTClaims{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	 }
	 token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	 tokenString , err := token.SignedString([]byte(secretKey))
	 if err != nil {
		return "", err
	}

	return tokenString, nil
}


func(s *authService) SendDailySummary() error {
	 if err := s.repo.SendDailySummary(); err != nil {
		 return err 
	 }
	 return nil 
}