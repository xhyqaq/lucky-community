package service_context

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/model"
)

const (
	errKey     = "err"
	msgKey     = "msg"
	dataKey    = "data"
	flashKey   = "flash"
	userKey    = "user"
	unreadKey  = "unread"
	versionKey = "version"
	nameKey    = "name"
	titleKey   = "title"
)

type BaseContext struct {
	Ctx     *gin.Context
	session sessions.Session
	path    string
	name    string
	title   string
}

func DataContext(ctx *gin.Context) *BaseContext {
	stx := &BaseContext{
		Ctx:     ctx,
		path:    "/",
		session: sessions.Default(ctx),
		name:    "lucky",
	}
	return stx
}

// Check 检查授权
func (c *BaseContext) Check() bool {
	user := c.Auth()
	if user == nil {
		return false
	} else {
		return user.ID > 0
	}
}

func (c *BaseContext) Refresh(users *model.Users) {
	c.SetAuth(users)
}

// SetAuth 设置授权
func (c *BaseContext) SetAuth(users *model.Users) {
	s, _ := json.Marshal(users)
	c.session.Set(userKey, string(s))
	_ = c.session.Save()
}

func (c *BaseContext) Auth() *model.Users {
	var user *model.Users
	str := c.session.Get(userKey)
	if str == nil {
		return user
	}
	if v, ok := str.(string); ok {
		_ = json.Unmarshal([]byte(v), &user)
	}
	return user
}

// Redirect 处理跳转
func (c *BaseContext) Redirect() {
	c.Ctx.Redirect(http.StatusFound, c.path)
}

// clear 清除闪存消息
func (c *BaseContext) clear() {
	c.session.Delete(errKey)
	c.session.Delete(msgKey)
	c.session.Delete(flashKey)
	_ = c.session.Save()
}

// Back 返回上一页
func (c *BaseContext) Back() *BaseContext {
	c.path = c.Ctx.Request.RequestURI
	return c
}

func (c *BaseContext) Referer(path ...string) {
	s := c.Ctx.Request.Referer()
	if s == "" {
		s = "/"
	} else if len(path) > 0 {
		s = path[0]
	}

	c.To(s)
	c.Redirect()
}

// To 设置跳转路径
func (c *BaseContext) To(to string) *BaseContext {
	c.path = to
	return c
}

// WithData 闪存消息
func (c *BaseContext) WithData(data interface{}) *BaseContext {
	r, _ := json.Marshal(data)
	c.session.Set(flashKey, string(r))
	_ = c.session.Save()
	return c
}

// ParseFlash 解析闪存数据
func (c *BaseContext) ParseFlash() map[string]interface{} {
	flashData := make(map[string]interface{})
	if str := c.session.Get(flashKey); str != nil {
		if v, ok := str.(string); ok {
			_ = json.Unmarshal([]byte(v), &flashData)
		}
	}
	return flashData
}

// WithError 错误信息跳转
func (c *BaseContext) WithError(err interface{}) *BaseContext {
	errStr := ""
	switch v := err.(type) {
	case error:
		errStr = v.Error()
	case string:
		errStr = v
	}
	c.session.Set(errKey, errStr)
	_ = c.session.Save()
	return c
}

// WithMsg 提示消息跳转
func (c *BaseContext) WithMsg(msg string) *BaseContext {
	c.session.Set(msgKey, msg)
	_ = c.session.Save()
	return c
}

// Forget 清除授权
func (c *BaseContext) Forget() {
	c.session.Delete(userKey)
	_ = c.session.Save()
}

// SetTitle 设置模版标题
func (c *BaseContext) SetTitle(title string) *BaseContext {
	c.title = title
	return c
}

// View 模版返回
func (c *BaseContext) View(tpl string, data interface{}) {
	obj := gin.H{
		versionKey: "1.0",
		errKey:     c.session.Get(errKey),
		msgKey:     c.session.Get(msgKey),
		userKey:    c.Auth(),
		dataKey:    data,
		flashKey:   c.ParseFlash(),
		nameKey:    c.name,
		titleKey:   c.title,
	}
	c.clear()
	c.Ctx.HTML(http.StatusOK, tpl, obj)
}

// Json 通用 JSON 响应
func (c *BaseContext) Json(data interface{}) {
	c.Ctx.JSON(http.StatusOK, data)
}

// MDFileJson markdown 上传图片响应
func (c *BaseContext) MDFileJson(ok int, msg, url string) {
	c.Json(gin.H{"success": ok, "message": msg, "url": url})
}
