package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salasasa/go-publisher/database/gorm"
	handler "github.com/salasasa/go-publisher/handler/gin"
	"github.com/salasasa/go-publisher/util"
)

func Init() {
	util.InitSlog("./output/go-publisher.log")
	gorm.ConnertPostDB("./conf", "db.yaml", util.YAML, "./output/")
}

func main() {
	Init()
	engine := gin.Default()
	engine.Static("/js", "./views/js")
	engine.Static("/css", "./views/css")
	engine.LoadHTMLGlob("./views/html/*")
	engine.StaticFile("/favicon.ico", "./views/img/dqq.png")

	engine.GET("/login", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "login.html", nil) })
	engine.GET("/regist", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "user_regist.html", nil) })
	engine.GET("/modify_pass", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "update_pass.html", nil) })
	engine.POST("/login/submit", handler.Login)
	engine.POST("/regist/submit", handler.RegistUser)
	engine.POST("/modify_pass/submit", handler.Auth, handler.UpdatePassword)
	engine.GET("/user", handler.GetUserInfo)
	engine.GET("/logout", handler.Logout)

	group := engine.Group("/news")
	group.GET("", handler.NewsList)
	group.GET("/issue", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "news_issue.html", nil) })
	group.POST("/issue/submit", handler.Auth, handler.PostNews)
	group.GET("/belong", handler.NewsBelong)
	group.GET("/:id", handler.GetNewsById)
	group.GET("/delete/:id", handler.Auth, handler.DeleteNews)
	group.POST("/update", handler.Auth, handler.UpdateNews)

	engine.GET("", func(ctx *gin.Context) { ctx.Redirect(http.StatusMovedPermanently, "news") }) //新闻列表页是默认的首页

	engine.Run("0.0.0.0:20001")
}
