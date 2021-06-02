package repository

import (
	"context"
	"fmt"

	"github.com/ismtabo/mapon-viewer/pkg/errors"
	"github.com/ismtabo/mapon-viewer/pkg/helper"
	"github.com/ismtabo/mapon-viewer/pkg/model"
	"github.com/ismtabo/mapon-viewer/pkg/repository/dao"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
)

// UserRepository implements user storage.
type UserRepository interface {
	// CreateUser creates an user in the application repository.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUser finds an user by email from the application repository.
	GetUser(ctx context.Context, email string) (*model.User, error)
}

type mongoUserRepository struct {
	collection helper.CollectionHelper
}

// NewMongoUserRepository creates an instance of UserRepository given a Mongo collection.
func NewMongoUserRepository(collection helper.CollectionHelper) UserRepository {
	return &mongoUserRepository{collection: collection}
}

// CreateUser creates an user in the given a Mongo collection.
func (r mongoUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	if exists, err := r.existsUser(ctx, user.Email); err != nil {
		return err
	} else if exists {
		return errors.NewConflictError(fmt.Sprintf("user with email '%s' already exists", user.Email))
	}
	userDAO := &dao.User{}
	copier.Copy(userDAO, user)
	if _, err := r.collection.InsertOne(ctx, userDAO); err != nil {
		return mapMongoError(err)
	}
	return nil
}

// GetUser finds an user by email in the given a Mongo collection.
func (r mongoUserRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email}
	result := r.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, mapMongoError(result.Err())
	}
	userDAO := &dao.User{}
	if err := result.Decode(userDAO); err != nil {
		return nil, err
	}
	user := &model.User{}
	if err := copier.Copy(user, userDAO); err != nil {
		return nil, err
	}
	return user, nil
}

func (r mongoUserRepository) existsUser(ctx context.Context, email string) (bool, error) {
	if _, err := r.GetUser(ctx, email); err != nil {
		if appErr, ok := err.(*errors.Error); !ok {
			return false, err
		} else if appErr.Code != "not_found" {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
