package services

import (
	"context"
	"errors"

	"github.com/lailiseptiandi/be-kominfo-test/internal/dtos/requests"
	"github.com/lailiseptiandi/be-kominfo-test/internal/repository"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/helper"
	"github.com/lailiseptiandi/be-kominfo-test/pkg/utils/jwt"
)

type AuthService interface {
	Login(ctx context.Context, req requests.UserRequest) (*string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *authService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(ctx context.Context, req requests.UserRequest) (*string, error) {

	us, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !helper.CheckPasswordHash(req.Password, us.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// generate token
	token, err := jwt.GenerateAccessToken(us)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	var userToken *string = &token

	return userToken, nil

}
