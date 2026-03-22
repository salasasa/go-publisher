package model

type User struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	Password string
}
