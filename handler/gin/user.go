package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/salasasa/go-publisher/database/gorm"
	"github.com/salasasa/go-publisher/handler/model"
	"github.com/salasasa/go-publisher/util"
)

func RegistUser(ctx *gin.Context) {
	user := &model.User{}

	if err := ctx.ShouldBind(user); err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}

	if err := gorm.RegistUser(user.Name, user.Password); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func Login(ctx *gin.Context) {
	user := &model.User{}
	ctx.ShouldBind(user)
	userData := gorm.GetUserByName(user.Name)
	if userData == nil {
		ctx.String(http.StatusBadRequest, "用户名不存在")
		return
	} else if userData.Password != user.Password {
		ctx.String(http.StatusBadRequest, "密码错误")
		return
	}

	ctx.SetCookie("uid", strconv.Itoa(userData.Id), 86400, "/", "192.168.10.110", false, true)
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("uid", "", -1, "/", "192.168.10.110", false, true)
}

func UpdatePassword(ctx *gin.Context) {
	req := &model.ModifyPassRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
		return
	}

	uid := GetUidFromCookie(ctx)
	if uid <= 0 {
		ctx.String(http.StatusBadRequest, "请先登录")
		return
	}

	if err := gorm.UpdateUserPassword(uid, req.NewPass, req.OldPass); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func GetUidFromCookie(ctx *gin.Context) int {
	for _, cookie := range ctx.Request.Cookies() {
		if cookie.Name == "uid" {
			if uid, err := strconv.Atoi(cookie.Value); err == nil {
				return uid
			}
		}
	}
	return 0
}
