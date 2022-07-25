package service

import (
	"context"

	"github.com/1ch0/go-restful/pkg/apiserver/infrastructure/datastore"
)

type UserService interface {
	Init(ctx context.Context) error
}

type userServiceImpl struct {
	Store datastore.DataStore `inject:"datastore"`
}

// NewUserService new User service
func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) Init(ctx context.Context) error {
	return nil
}
