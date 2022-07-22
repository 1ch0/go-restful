package service

import (
	"context"
)

type UserService interface {
	Init(ctx context.Context) error
}

type userServiceImpl struct {
}

// NewUserService new User service
func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) Init(ctx context.Context) error {
	return nil
}
