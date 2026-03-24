package gorm_test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"testing"

	"github.com/salasasa/go-publisher/database/gorm"
	"github.com/salasasa/go-publisher/util"
)

func init() {
	util.InitSlog("../../output/go-publisher.log")
	gorm.ConnertPostDB("../../conf", "db.yaml", util.YAML, "../../output/")
}

func hash(pass string) string {
	hasher := md5.New()
	hasher.Write([]byte(pass))
	digest := hasher.Sum(nil)
	return hex.EncodeToString(digest)
}

var (
	uid = 5
)

func TestRegistUser(t *testing.T) {
	slog.Info("try")
	err := gorm.RegistUser("ArcE", hash("123456"))
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("注册成功\n")
	}

	err = gorm.RegistUser("ArcE", hash("123456"))
	if err != nil {
		fmt.Printf("注册失败: %s\n", err)
	} else {
		fmt.Println("重复注册成功！")
		t.Fail()
	}
}

func TestLogOffUser(t *testing.T) {
	if err := gorm.LogOffUser(1); err != nil {
		slog.Error(err.Error())
		t.Fail()
	}
}
