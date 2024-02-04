package services

import (
	"errors"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserService struct {
}

// get user information
func (*UserService) GetUserById(id int) *model.Users {

	user := userDao.QueryUser(&model.Users{ID: id})
	user.Avatar = utils.BuildFileUrl(user.Avatar)
	return user
}

// update user information
func (*UserService) UpdateUser(user *model.Users) {

	userDao.UpdateUser(user)

}

func (*UserService) ListByIdsSelectEmail(id ...int) []string {
	return userDao.ListByIds(id...)
}

func (s *UserService) ListByIdsSelectIdNameMap(ids []int) map[int]string {

	m := make(map[int]string)
	users := userDao.ListByIdsSelectIdName(ids)
	for i := range users {
		user := users[i]
		m[user.ID] = user.Name
	}
	return m
}

func Login(account, pswd string) (*model.Users, error) {

	user := userDao.QueryUser(&model.Users{Account: account, Password: pswd})
	if user.ID == 0 {
		return nil, errors.New("登录失败！账号或密码错误")
	}

	user.Avatar = utils.BuildFileUrl(user.Avatar)
	return user, nil
}

func Register(account, pswd, name string, inviteCode uint16) error {

	if err := utils.NotBlank(account, pswd, name, inviteCode); err != nil {
		return err
	}

	// query codeDao
	if !codeDao.Exist(inviteCode) {
		return errors.New("验证码不存在")
	}

	// 查询账户
	user := userDao.QueryUser(&model.Users{Account: account})
	if user.ID > 0 {
		return errors.New("账户已存在,换一个吧")
	}

	// 保存用户
	userDao.CreateUser(account, name, pswd, inviteCode)
	// 修改code状态
	SetState(inviteCode)

	return nil
}
