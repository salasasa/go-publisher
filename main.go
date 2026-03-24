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

	engine.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	engine.GET("/regist", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "user_regist.html", nil)
	})
	engine.GET("/modify_pass", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "update_pass.html", nil)
	})

	engine.POST("/login/submit", handler.Login)
	engine.POST("/regist/submit", handler.RegistUser)
	engine.POST("/modify_pass/submit", handler.UpdatePassword)
	engine.GET("/logout", handler.Logout)

	engine.Run("0.0.0.0:20001")

}
