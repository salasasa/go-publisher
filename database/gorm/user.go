package gorm

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/salasasa/go-publisher/database/model"
	"gorm.io/gorm"
)

func RegistUser(name, password string) error {
	user := model.User{
		Name:     name,
		Password: password,
	}
	if err := GoPublisherDB.Create(&user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return fmt.Errorf("用户名[%s]已存在", name)
			}
		}
		slog.Error("注册用户失败", "name", name, "error", err)
		return fmt.Errorf("注册用户[%s]失败: %w", name, err)
	}
	return nil
}
func LogOffUser(uid int) error {
	user := &model.User{
		Id: uid,
	}

	tx := GoPublisherDB.
		// Session(&gorm.Session{DryRun: true}). // 测试阶段，先不要真删
		Delete(user)
	if tx.Error != nil {
		slog.Error("注销用户失败", "uid", uid, "error", tx.Error)
		return errors.New("用户注销失败，请稍后重试")
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("用户注销失败，uid %d不存在", uid)
	}
	return nil
}

func GetUserById(uid int) *model.User {
	user := &model.User{}
	tx := GoPublisherDB.Select("*").Where("id = ?", uid).First(user)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("GetUserById failed", "uid", uid, "error", tx.Error)
		}
		return nil
	}
	return user
}

func GetUserByName(name string) *model.User {
	user := &model.User{}
	tx := GoPublisherDB.Select("*").Where("name = ?", name).First(user)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("GetUserByName failed", "name", name, "error", tx.Error)
		}
		return nil
	}

	return user
}

func UpdateUserName(uid int, name string) error {
	tx := GoPublisherDB.Model(&model.User{}).Where("id = ?", uid).Update("name", name)
	if tx.Error != nil {
		slog.Error("UpdateUserName failed", "uid", uid, "new name", name, "error", tx.Error)
		return errors.New("用户名修改失败，请稍后重试")
	}

	if tx.RowsAffected <= 0 {
		return fmt.Errorf("用户id[%d]不存在", uid)
	}

	return nil
}

func UpdateUserPassword(uid int, newPass, oldPass string) error {
	tx := GoPublisherDB.Model(&model.User{}).Where("id = ?", uid).Where("password = ?", oldPass).Update("password", newPass)
	if tx.Error != nil {
		slog.Error("UpdatePassword failed", "uid", uid, "error", tx.Error)
		return errors.New("密码修改失败，请稍后重试")
	}

	if tx.RowsAffected <= 0 {
		return errors.New("旧密码不对")
	}

	return nil
}
