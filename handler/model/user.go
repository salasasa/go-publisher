package model

type User struct {
	Name     string `form:"name" binding:"required,gte=2"`
	Password string `form:"pass" binding:"required,len=32"`
}

type ModifyPassRequest struct {
	OldPass string `form:"old_pass" binding:"required,len=32" validate:"required,len=32"`
	NewPass string `form:"new_pass" binding:"required,len=32" validate:"required,len=32"`
}
