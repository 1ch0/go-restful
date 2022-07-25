package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/1ch0/go-restful/pkg/apiserver/domain/model"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/bcode"
	"github.com/1ch0/go-restful/pkg/apiserver/utils/log"
	"github.com/1ch0/go-restful/pkg/utils"
	"golang.org/x/crypto/bcrypt"

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
	admin := model.DefaultAdminUserName
	//pwd := func() string {
	//	p := utils.RandomString(8)
	//	p += strconv.Itoa(rand.Intn(9))                 // #nosec
	//	r := append([]rune(p), 'a'+rune(rand.Intn(26))) // #nosec
	//	rand.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })
	//	p = string(r)
	//	return p
	//}()
	pwd := utils.SetPasswd(8)
	//user := &model.User{Name: admin}
	if err := u.Store.Get(ctx, &model.User{
		Name: admin,
	}); err != nil {
		if errors.Is(err, datastore.ErrRecordNotExist) {
			encrypted, err := GeneratePasswordHash(pwd)
			if err != nil {
				return err
			}
			if err := u.Store.Add(ctx, &model.User{
				Name:     admin,
				Password: encrypted,
			}); err != nil {
				fmt.Println("------------------")
				return err
			}
			// print default password of admin user in log
			log.Logger.Infof("initialized admin username and password: admin / %s", pwd)
		} else {
			return err
		}
	}
	//log.Logger.Infof("admin user is exist, pwd is: %s", pwd)
	return nil
}

func GeneratePasswordHash(s string) (string, error) {
	if s == "" {
		return "", bcode.ErrUserInvalidPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
