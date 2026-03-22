package gorm

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/salasasa/go-publisher/database/model"
)

func RegistUser(name, password string) (int, error) {
	user := model.User{
		Name:     name,
		Password: password,
	}

	if err := GoPublisherDB.Create(&user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return 0, fmt.Errorf("用户名[%s]已存在", name)
			}
		}
		slog.Error("注册用户失败", "name", name, "error", err)
		return 0, fmt.Errorf("注册用户[%s]失败: %w", name, err)
	}
	return user.Id, nil
}
