package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	xt "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

type registerForm struct {
	Code     int    `binding:"required" form:"code" msg:"code不能为空" `
	Account  string `binding:"required,email" form:"account" msg:"邮箱格式不正确"`
	Name     string `binding:"required" form:"name" msg:"用户名不能为空"`
	Password string `binding:"required" form:"password" msg:"密码不能为空"`
}

func InitLoginRegisterRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	group.POST("/login", Login)
	group.POST("/register", Register)
}

func Login(c *gin.Context) {
	var login model.LoginForm
	if err := c.ShouldBindJSON(&login); err != nil {
		result.Err(utils.GetValidateErr(login, err)).Json(c)
		return
	}
	loginLog := model.LoginLogs{
		Account:   login.Account,
		Browser:   c.Request.UserAgent(),
		Equipment: c.GetHeader("Sec-Ch-Ua-Platform"),
		Ip:        utils.GetClientIP(c.Request),
		CreatedAt: xt.Now(),
	}
	var logS services.LogServices
	user, err := services.Login(login)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	token, err := middleware.GenerateToken(user.ID, user.Name)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	c.SetCookie(middleware.AUTHORIZATION, token, 3600, "/", c.Request.Host, true, true)
	loginLog.State = "登录成功"
	logS.InsertLoginLog(loginLog)
	result.OkWithMsg(map[string]string{"token": token}, "登录成功").Json(c)
}

func Register(c *gin.Context) {
	var form registerForm

	err := c.ShouldBindJSON(&form)
	loginLog := model.LoginLogs{
		Account:   form.Account,
		Browser:   c.Request.UserAgent(),
		Equipment: c.GetHeader("Sec-Ch-Ua-Platform"),
		Ip:        utils.GetClientIP(c.Request),
		CreatedAt: xt.Now(),
	}
	var logS services.LogServices
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(utils.GetValidateErr(form, err)).Json(c)
		return
	}

	id, err := services.Register(form.Account, form.Password, form.Name, form.Code)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	loginLog.State = "注册成功"
	logS.InsertLoginLog(loginLog)
	token, err := middleware.GenerateToken(id, form.Name)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	c.SetCookie(middleware.AUTHORIZATION, token, 3600, "/", c.Request.Host, true, true)
	result.OkWithMsg(map[string]string{"token": token}, "注册成功").Json(c)
}
