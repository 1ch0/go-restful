package service

import (
	"context"
	"fmt"

	"github.com/1ch0/go-restful/pkg/apiserver/config"
)

var needInitData []DataInit

func InitServiceBean(c config.Config) []interface{} {
	//rbacService := NewRBACService()
	userService := NewUserService()
	authenticationService := NewAuthenticationService()
	needInitData = []DataInit{userService}
	return []interface{}{userService, authenticationService}
}

type DataInit interface {
	Init(ctx context.Context) error
}

func InitData(ctx context.Context) error {
	for _, init := range needInitData {
		if err := init.Init(ctx); err != nil {
			return fmt.Errorf("database init failure %w", err)
		}
	}
	return nil
}
