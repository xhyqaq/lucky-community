package frontend

import (
	"fmt"
	"strconv"

	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/constants"
	"xhyovo.cn/community/server/service/event"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitCommentRouters(g *gin.Engine) {
	group := g.Group("/community/comments")
	group.GET("/byArticleId/:articleId", listCommentsByArticleId)
	group.GET("/byRootId/:rootId", listCommentsByRootId)
	group.GET("/allCommentsByArticleId/:articleId", listAllCommentsByArticleId)
	group.GET("/adaptions", adaptions)
	group.GET("/byArticleId", listCommentsByArticleIdNoTree)
	group.GET("/latest", listLatestComments)
	group.GET("/summary/:businessId", getCommentSummary)
	group.Use(middleware.OperLogger())
	group.POST("/comment", comment)
	group.DELETE("/:id", deleteComment)
	group.POST("/adoption", adoption)

}

// 获取采纳评论
func adaptions(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Query("articleId"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户id: %d 获取采纳评论解析参数失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	p, limit := page.GetPage(ctx)
	commentsService := services.NewCommentService(ctx)

	comments, count := commentsService.ListAdoptionsByArticleId(articleId, p, limit)
	result.Page(comments, count, nil).Json(ctx)
}

// 发布评论
func comment(ctx *gin.Context) {
	var comment model.Comments
	userId := middleware.GetUserId(ctx)
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		msg := utils.GetValidateErr(comment, err)
		log.Warnf("用户id: %d 发布评论失败,err: %s", userId, msg)
		result.Err(msg).Json(ctx)
		return
	}
	comment.FromUserId = userId

	commentsService := services.NewCommentService(ctx)
	err := commentsService.Comment(&comment)
	msg := "评论成功"
	if err != nil {
		log.Warnf("用户id: %d 保存评论失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, msg).Json(ctx)
}

// 删除评论
func deleteComment(ctx *gin.Context) {
	commentId := ctx.Param("id")

	userId := middleware.GetUserId(ctx)
	if commentId == "" {
		log.Warnf("用户id: %d 删除评论失败,err: %s", userId, "评论id为空")
		result.Err("删除评论id不能为空").Json(ctx)
		return
	}
	commentIdInt, _ := strconv.Atoi(commentId)
	var commentsService services.CommentsService
	if !commentsService.DeleteComment(commentIdInt, userId) {
		log.Warnf("用户id: %d 删除评论失败", userId)
		result.Err("删除失败").Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 返回文章下的评论(文章页面展示)
func listCommentsByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("articleId"))
	p, limit := page.GetPage(ctx)
	currentUserId := middleware.GetUserId(ctx)

	if err != nil {
		log.Warnf("用户id: %d 获取文章下的评论失败,err: %s", currentUserId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByArticleID(p, limit, articleId, currentUserId)
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(page.New(comments, count), "").Json(ctx)
}

// 查询根评论下的评论
func listCommentsByRootId(ctx *gin.Context) {
	rootId, _ := strconv.Atoi(ctx.Param("rootId"))
	p, limit := page.GetPage(ctx)
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByRootID(p, limit, rootId)
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(page.New(comments, count), "").Json(ctx)

}

// 查询用户文章下的所有评论，文章id为空则查询所有(管理端)
func listAllCommentsByArticleId(ctx *gin.Context) {

	p, limit := page.GetPage(ctx)

	userId := middleware.GetUserId(ctx)
	var commentsService services.CommentsService
	comments, count := commentsService.GetAllCommentsByArticleID(p, limit, userId, 0, 0)
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(page.New(comments, count), "").Json(ctx)
}

// 采纳评论
func adoption(ctx *gin.Context) {
	var adoption model.QaAdoptions
	userId := middleware.GetUserId(ctx)
	if err := ctx.ShouldBindJSON(&adoption); err != nil {
		msg := utils.GetValidateErr(adoption, err)
		log.Warnf("用户id: %d 采纳评论参数解析失败,err :%s", userId, msg)
		result.Err(msg).Json(ctx)
		return
	}
	var cS services.CommentsService
	commentId := adoption.CommentId
	comment := cS.GetById(commentId)
	articleId := comment.BusinessId
	adoption.ArticleId = articleId
	var msg string
	// 采纳权限，文章得是本人,评论得存在
	var aS services.ArticleService
	article := aS.GetById(articleId)
	state := article.State
	if article.UserId != userId {
		msg = fmt.Sprintf("用户id: %d 采纳评论无权限,文章id: %d", userId, articleId)
		log.Warnln(msg)
		result.Err("只有发布者运行采纳").Json(ctx)
		return
	}
	if state != constant.PrivateQuestion {
		if !(state == constants.Pending || state == constants.Resolved) {
			result.Err("该文章不是 QA 分类,无法进行采纳").Json(ctx)
			return
		}
	}

	if !aS.Auth(userId, articleId) {
		msg = fmt.Sprintf("用户id: %d 采纳评论无权限,文章id: %d", userId, articleId)
		log.Warnln(msg)
		result.Err(msg).Json(ctx)
		return
	}

	if comment.ID == 0 {
		msg = fmt.Sprintf("用户id: %d 采纳评论对应的评论不存在,评论id: %d", userId, commentId)
		log.Warnln(msg)
		result.Err(msg).Json(ctx)
		return
	}
	var adptionS services.QAAdoption
	msg = "取消采纳"
	if adptionS.Adopt(articleId, commentId) {
		var suS services.SubscriptionService
		suS.Send(event.Adoption, constant.NOTICE, userId, comment.FromUserId, services.SubscribeData{CommentId: commentId, ArticleId: articleId, UserId: userId, CurrentBusinessId: articleId})
		msg = "已采纳"
	}
	if state != constant.PrivateQuestion {
		// 采纳了,但是状态为未解决,则改为已解决
		if adptionS.QAAdoptState(articleId) && state == constants.Pending {
			aS.UpdateState(articleId, constants.Resolved)
		} else if !adptionS.QAAdoptState(articleId) && state == constants.Resolved {
			aS.UpdateState(articleId, constants.Pending)
		}
	} else {
		msg = "已采纳,当前版本私密提问采纳后无法变更为解决,请等待"
	}

	result.OkWithMsg(nil, msg).Json(ctx)
}

// 返回文章下的所有评论，非树形结构
func listCommentsByArticleIdNoTree(ctx *gin.Context) {
	businessId, err := strconv.Atoi(ctx.Query("businessId"))
	if err != nil {
		result.Err("获取文章下的所有评论,文章 id 解析失败, err:" + err.Error()).Json(ctx)
		return
	}
	tenantId, err := strconv.Atoi(ctx.Query("tenantId"))
	if err != nil {
		log.Warnf("用户id: %d 查询用户文章下的所有评论失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("查询对应模块id不可为空").Json(ctx)
		return
	}

	currentUserId := middleware.GetUserId(ctx)
	cS := services.NewCommentService(ctx)
	comments := cS.ListCommentsByArticleIdNoTree(businessId, tenantId, currentUserId)
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(comments, "").Json(ctx)
}

// listLatestComments 获取最新10条评论
func listLatestComments(ctx *gin.Context) {
	var c services.CommentsService
	comments, count := c.ListLatestComments()
	result.Page(comments, count, nil).Json(ctx)
}

// 获取评论总结
func getCommentSummary(ctx *gin.Context) {
	businessId, err := strconv.Atoi(ctx.Param("businessId"))
	if err != nil {
		log.Warnf("获取评论总结参数解析失败,err: %s", err.Error())
		result.Err("业务ID参数错误").Json(ctx)
		return
	}
	
	tenantId, err := strconv.Atoi(ctx.Query("tenantId"))
	if err != nil {
		log.Warnf("获取评论总结租户ID解析失败,err: %s", err.Error())
		result.Err("租户ID参数错误").Json(ctx)
		return
	}
	
	summaryService := services.NewCommentSummaryService(ctx)
	summary, err := summaryService.GetSummary(businessId, tenantId)
	if err != nil {
		log.Warnf("获取评论总结失败,err: %s", err.Error())
		result.Err("获取总结失败").Json(ctx)
		return
	}
	
	// 如果没有评论，返回空结果
	if summary == nil {
		result.Ok(nil, "暂无评论总结").Json(ctx)
		return
	}
	
	result.Ok(summary, "获取成功").Json(ctx)
}
