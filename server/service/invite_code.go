package services

import (
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/dao"
)

func GenerateCode() uint16 {

	var flag bool
	var code uint16
	// generate code
	for flag {
		code = utils.GenerateCode(8)
		flag = dao.InviteCode.Exist(code)
	}

	return code
}

func DestroyCode(code int) error {

	return dao.InviteCode.Del(code)
}

func SetState(code uint16) error {
	return dao.InviteCode.SetState(code)
}