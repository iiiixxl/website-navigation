package model

import "time"

type User struct {
	ID           int       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Account      string    `json:"account" gorm:"column:account;uniqueIndex:uk_account;size:50"`
	Password     string    `json:"-" gorm:"column:password;size:255;not null"`
	Salt         string    `json:"-" gorm:"column:salt;size:256"`
	Name         string    `json:"name" gorm:"column:name;size:50;not null"`
	EnName       string    `json:"en_name" gorm:"column:en_name;size:100"`
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time"`
	CreateUserID *int      `json:"create_user_id" gorm:"column:create_user_id"`
	IsDeleted    bool      `json:"is_deleted" gorm:"column:is_deleted;default:0;not null"`
}

func (User) TableName() string {
	return "t_user"
}

type CreateUserRequest struct {
	Account  string `json:"account" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Name     string `json:"name" binding:"required" example:"管理员"`
	EnName   string `json:"en_name" example:"Admin"`
}

type UpdateUserRequest struct {
	ID       int    `json:"id" binding:"required" example:"1"`
	Password string `json:"password" example:"123456"`
	Name     string `json:"name" example:"管理员"`
	EnName   string `json:"en_name" example:"Admin"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
