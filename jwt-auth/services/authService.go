package services

import (
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/domain"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/dto"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, error)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (a DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	user, err := a.repo.ValidateUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := user.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func NewAuthService(repo domain.AuthRepository) DefaultAuthService {
	return DefaultAuthService{repo: repo}
}
