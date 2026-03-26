package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salasasa/go-publisher/util"
)

const (
	UID_IN_TOKEN = "uid"
	UID_IN_CTX   = "uid"
	COOKIE_NAME  = "jwt"
)

var (
	KeyConfig = util.InitViper("conf", "jwt", util.YAML)
)

func Auth(ctx *gin.Context) {
	LoginUid := GetLoginUid(ctx)
	if LoginUid <= 0 {
		ctx.String(http.StatusUnauthorized, "请先登录")
		ctx.Redirect(http.StatusTemporaryRedirect, "/login") //重定向到登录页面
		ctx.Abort()
		return
	}
	ctx.Set(UID_IN_CTX, LoginUid)
}

func GetLoginUid(ctx *gin.Context) int {
	cookies := ctx.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == COOKIE_NAME {
			return GetUidFromJwt(cookie.Value)
		}
	}

	return 0
}

func GetUidFromJwt(token string) int {
	_, jwtPayload, err := util.VerifyJwt(token, KeyConfig.GetString("secret"))
	if err != nil {
		slog.Error("VerifyJwt failed", "error", err)
		return 0
	}
	userData := jwtPayload.UserDefined
	if userData == nil {
		slog.Error("userData is nil")
		return 0
	}
	uid, ok := userData[UID_IN_TOKEN].(float64)
	if !ok {
		slog.Error("uid is not int")
		return 0
	}
	return int(uid)
}
