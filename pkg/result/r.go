package result

import "github.com/gin-gonic/gin"

type R struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Ok(data any, msg string) *R {
	return &R{Code: 200, Data: data, Msg: msg}
}
func Err(msg string) *R {
	return &R{Code: 500, Data: nil, Msg: msg}
}

func (r *R) Json(c *gin.Context) {
	c.JSON(r.Code, r)
}

func (r *R) Xml(c *gin.Context) {
	c.XML(r.Code, r)
}
