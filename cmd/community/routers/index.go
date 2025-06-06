package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var courseService services.CourseService

func InitIndexRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	// 社区首页
	group.GET("/user/count", getUserCount)
	group.GET("/rate/page", pageRate)
	group.GET("/labels/", getKnowedgeLabels)
	group.GET("/index/courses", GetHomepageCourses)
}

func pageRate(ctx *gin.Context) {

	p, limit := page.GetPage(ctx)
	var noteService services.RateService
	state, notes := noteService.Page(p, limit)

	result.Ok(map[string]interface{}{
		"state": state,
		"data":  notes,
	}, "").Json(ctx)
	return
}

func getUserCount(ctx *gin.Context) {
	var count int64
	model.User().Count(&count)
	result.Ok(count, "").Json(ctx)
}

func getKnowedgeLabels(ctx *gin.Context) {

	labels := []string{"javase", "juc", "jvm", "mysql", "redis", "mq", "多线程", "反射", "字节码", "设计模式", "spring", "springmvc", "mybatis", "springboot", "dubbo", "分布式", "微服务", "zookeeper", "计算机网络", "操作系统"}
	result.Ok(labels, "").Json(ctx)
}

// GetHomepageCourses 获取首页所有课程信息（包括详情和章节列表）
func GetHomepageCourses(ctx *gin.Context) {
	courses, err := courseService.GetAllCoursesWithDetails()
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Ok(courses, "").Json(ctx)
}
