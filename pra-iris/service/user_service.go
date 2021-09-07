package service

import (
	"errors"
	"github.com/go-xorm/xorm"
	"pra-iris/m/res"
	"pra-iris/models"
	"time"
)

type UserService interface {
	Add(username string, mobile string, password string, cityName string) res.Res
}

type userService struct {
	Engine *xorm.Engine
}

func NewUserService(db *xorm.Engine) UserService {
	return &userService{
		Engine: db,
	}
}

func (us *userService) Add(username string, mobile string, password string, cityName string) res.Res {
	if username == "" {
		return res.MyRes(
			-1,
			errors.New("username is not allowed null"),
			errors.New("用户名不能为空"),
		)
	}

	var user models.User
	has, err := us.Engine.Where(" user_name = ? ", username).Get(&user)

	if has {
		return res.MyRes(
			-3,
			//errors.New(user.UserName),
			errors.New("user exists"),
			errors.New("用户已存在：" + user.UserName),
			)
	}

	if password == "" {
		return res.MyRes(
			-2,
			errors.New("password is not allowed null"),
			//errors.New(err.Error()),
			errors.New("密码不能为空"),
		)
	}

	user = models.User{
		UserName: username,
		RegisterTime: time.Now(),
		Mobile: mobile,
		Pwd: password,
		CityName: cityName,
	}

	_, err = us.Engine.Insert(user)

	return res.MyRes(1, err, err)
}