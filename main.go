package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/salasasa/go-publisher/database/gorm"
	handler "github.com/salasasa/go-publisher/handler/gin"
	"github.com/salasasa/go-publisher/util"
)

func Init() {
	util.InitSlog("./output/go-publisher.log")
	gorm.ConnertPostDB("./conf", "db.yaml", util.YAML, "./output/")

	crontab := cron.New()
	crontab.AddFunc("*/21 * * * *", gorm.PingPostDB) // 分，时，日，月，星期。每隔21分钟ping一次数据库
	crontab.Start()
}

func ListenTermSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15。Ctrl+C对应SIGINT信号
	sig := <-c                                        //阻塞，直到信号的到来
	slog.Info("receive term signal " + sig.String() + ", going to exit")
	gorm.ClosePostDB() //关闭数据库连接
	os.Exit(0)         //进程退出
}

func main() {
	Init()
	go ListenTermSignal()

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
