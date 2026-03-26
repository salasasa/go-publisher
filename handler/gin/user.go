package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/salasasa/go-publisher/database/gorm"
	"github.com/salasasa/go-publisher/handler/model"
	"github.com/salasasa/go-publisher/util"
)

const (
	COOKIE_LIFE = 7 * 86400
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

	jwtHeader := util.DefautJwtHeader
	jwtpayload := util.JwtPayload{ //payload以明文形式编码在token中，server用自己的密钥可以校验该信息是否被篡改过
		Issue:       "news",
		IssueAt:     time.Now().Unix(),                                //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(COOKIE_LIFE * time.Second).Unix(), //7天后过期
		UserDefined: map[string]any{UID_IN_TOKEN: userData.Id},        //用户自定义字段。如果token里包含敏感信息，请结合https使用
	}
	if token, err := util.GenJWT(jwtHeader, jwtpayload, KeyConfig.GetString("secret")); err != nil {
		slog.Error("生成JWT失败", "error", err)
		ctx.String(http.StatusInternalServerError, "token生成失败")
		return
	} else {
		ctx.SetCookie(COOKIE_NAME, token, COOKIE_LIFE, "/", "192.168.10.110", false, true)
	}

	// ctx.SetCookie("uid", strconv.Itoa(userData.Id), 86400, "/", "192.168.10.110", false, true)
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie(COOKIE_NAME, "", -1, "/", "192.168.10.110", false, true)
}

func UpdatePassword(ctx *gin.Context) {
	uid, ok := ctx.Value(UID_IN_CTX).(int)
	if !ok {
		ctx.String(http.StatusBadRequest, "请先登录")
		return
	}

	req := &model.ModifyPassRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		ctx.String(http.StatusBadRequest, util.BindErrMsg(err))
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
	// for _, cookie := range ctx.Request.Cookies() {
	// 	if cookie.Name == COOKIE_NAME {
	// 		return GetUidFromJwt(cookie.Value)
	// 	}
	// }
	return 0
}
