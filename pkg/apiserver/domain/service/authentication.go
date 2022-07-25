package service

import (
	"context"
	"fmt"

	"github.com/1ch0/go-restful/pkg/apiserver/infrastructure/datastore"

	apisv1 "github.com/1ch0/go-restful/pkg/apiserver/interface/api/dto/v1"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/bcode"
)

const (
	GrantTypeAccess  = "access"
	GrantTypeRefresh = "refresh"
)

var signeKey = ""

type AuthenticationService interface {
	Login(ctx context.Context, loginReq apisv1.LoginRequest) (*apisv1.LoginResponse, error)
}

type authenticationServiceImpl struct {
	UserService UserService         `inject:""`
	Store       datastore.DataStore `inject:"datastore"`
}

func NewAuthenticationService() AuthenticationService {
	return &authenticationServiceImpl{}
}

type authHanler interface {
	login(ctx context.Context) (*apisv1.UserBase, error)
}

type localHandlerImpl struct {
	ds          datastore.DataStore
	userService UserService
	username    string
	password    string
}

func (a *authenticationServiceImpl) newLocalHandler(req apisv1.LoginRequest) (*localHandlerImpl, error) {
	if req.Username == "" || req.Password == "" {
		return nil, bcode.ErrInvalidLoginRequest
	}
	return &localHandlerImpl{
		//userService: a.UserService,
		username: req.Username,
		password: req.Password,
	}, nil
}

func (a *authenticationServiceImpl) Login(ctx context.Context, loginReq apisv1.LoginRequest) (*apisv1.LoginResponse, error) {
	var handler authHanler
	var err error

	handler, err = a.newLocalHandler(loginReq)
	if err != nil {
		return nil, err
	}
	userBase, err := handler.login(ctx)
	if err != nil {
		return nil, err
	}
	if userBase.Disabled {
		return nil, fmt.Errorf("TBD")
	}
	// TODO(@1CH0)
	return &apisv1.LoginResponse{
		User:         userBase,
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func (l *localHandlerImpl) login(ctx context.Context) (*apisv1.UserBase, error) {
	// TODO(@1CH0)
	//user, err := l.userService.GetUser(ctx, l.username)
	//if err != nil {
	//	if errors.Is(err, datastore.ErrRecordNotExist) {
	//		return nil, bcode.ErrUsernameNotExist
	//	}
	//	return nil, err
	//}
	//if err := compareHashWithPassword(user.Password, l.password); err != nil {
	//	return nil, err
	//}
	//if err := l.userService.UpdateUserLoginTime(ctx, user); err != nil {
	//	return nil, err
	//}
	return &apisv1.UserBase{
		Name:  "TBD",
		Email: "TBD",
	}, nil
}
