package auth

import (
	"errors"
	"time"
)


type AuthService interface {
     CreateUser(name, email , timezone string) error 

	 SendDailySummary() error 
}

type authService struct {
	 repo AuthRepo 
}

func NewAuthService (repo AuthRepo) AuthService {
	 return &authService{repo: repo}
}

func(s *authService) CreateUser(name, email , timezone string) error {
	 if name=="" {
		 return errors.New("name is required")
	 }
	 if email=="" {
		 return errors.New("email is required")
	 }
	 if timezone=="" {
		 return errors.New("timezone is required")
	 }
	 user:=&User{
		Name:name,
		Email: email,
		Timezone: timezone,
		CreatedAt: time.Now().UTC(),
	 }

	 err:=s.repo.Create(user)
	 if err != nil {
		 return err 
	 }
	 return  nil 

}




func(s *authService) SendDailySummary() error {
	 if err := s.repo.SendDailySummary(); err != nil {
		 return err 
	 }
	 return nil 
}