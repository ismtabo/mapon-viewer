package service

import (
	"context"

	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/model"
	"github.com/ismtabo/mapon-viewer/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements User handling use case.
type UserService interface {
	// ValidateUserPassword validates user authentication.
	ValidateUserPassword(ctx context.Context, email, password string) error
	// CreateUser creates a new user in the application.
	CreateUser(ctx context.Context, email, password string) error
}

type userService struct {
	repository repository.UserRepository
}

// NewUserService creates an new UserService.
func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository: repository}
}

// ValidateUserPassword validates given password for user identified by email.
func (s *userService) ValidateUserPassword(ctx context.Context, email, password string) error {
	user, err := s.repository.GetUser(ctx, email)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			if appErr.Code == errors.NotFoundErrorCode {
				return errors.NewUnauthorizedError().WithWrap(err)
			}
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.NewUnauthorizedError().WithWrap(err)
	}
	return nil
}

// CreateUser creates an application user with given email and password.
func (s *userService) CreateUser(ctx context.Context, email, password string) error {
	password, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}
	user := &model.User{Email: email, Password: password}
	return s.repository.CreateUser(ctx, user)
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", errors.NewInternalServerError(err)
	}
	return string(hash), nil
}
