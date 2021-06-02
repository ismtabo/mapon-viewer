package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/model"
	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/ismtabo/mapon-viewer/pkg/test/mocks/repository"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserServiceValidateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := repository.NewMockUserRepository(ctrl)
	service := service.NewUserService(repository)
	t.Run("Test UserService validate user password", func(t *testing.T) {
		pwd := "password"
		password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			t.Error(err)
		}
		user := &model.User{Email: "email", Password: string(password)}
		repository.EXPECT().GetUser(gomock.Any(), "email").Return(user, nil)
		err = service.ValidateUserPassword(context.Background(), "email", pwd)
		assert.NoError(t, err)
	})
	t.Run("Test UserService validate user password. User not found", func(t *testing.T) {
		expectedErr := errors.NewNotFoundError()
		repository.EXPECT().GetUser(gomock.Any(), "email").Return(nil, expectedErr)
		err := service.ValidateUserPassword(context.Background(), "email", "password")
		assert.ErrorIs(t, expectedErr, err)
	})
	t.Run("Test UserService validate user password. Internal server error", func(t *testing.T) {
		expectedErr := errors.NewInternalServerError(nil)
		repository.EXPECT().GetUser(gomock.Any(), "email").Return(nil, expectedErr)
		err := service.ValidateUserPassword(context.Background(), "email", "password")
		assert.ErrorIs(t, expectedErr, err)
	})
	t.Run("Test UserService validate user password. Incorrect password", func(t *testing.T) {
		password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
		if err != nil {
			t.Error(err)
		}
		user := &model.User{Email: "email", Password: string(password)}
		repository.EXPECT().GetUser(gomock.Any(), "email").Return(user, nil)
		err = service.ValidateUserPassword(context.Background(), "email", "incorrect")
		assert.Error(t, err)
	})
}

func TestUserServiceCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := repository.NewMockUserRepository(ctrl)
	service := service.NewUserService(repository)
	t.Run("Test UserService create user", func(t *testing.T) {
		pwd := "password"
		password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
		if err != nil {
			t.Error(err)
		}
		user := &model.User{Email: "email", Password: string(password)}
		repository.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(user)).Return(nil)
		err = service.CreateUser(context.Background(), "email", "password")
		assert.NoError(t, err)
	})
	t.Run("Test UserService create user. Conflict error", func(t *testing.T) {
		expectedErr := errors.NewConflictError("")
		repository.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).Return(expectedErr)
		err := service.CreateUser(context.Background(), "email", "password")
		assert.ErrorIs(t, expectedErr, err)
	})
	t.Run("Test UserService create user. Internal server error", func(t *testing.T) {
		expectedErr := errors.NewInternalServerError(nil)
		repository.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(&model.User{})).Return(expectedErr)
		err := service.CreateUser(context.Background(), "email", "password")
		assert.ErrorIs(t, expectedErr, err)
	})
}
