package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	UserRoleDefault = "user"
	UserRoleAdmin   = "admin"
)

type User struct {
	ID               primitive.ObjectID `json:"id"`
	Login            string             `json:"login"`
	Password         string             `json:"password"`
	Name             string             `json:"name"`
	Role             string             `json:"role"`
	TelegramUsername string             `json:"telegram_username"`
	TelegramChatID   int64              `json:"telegram_chat_id"`
	AuthCode         string             `json:"auth_code"`
	AuthCodeTime     int64              `json:"auth_code_time"`
}

type UserInfo struct {
	ID                 primitive.ObjectID `json:"id"`
	Name               string             `json:"name"`
	TelegramUsername   string             `json:"telegram_username"`
	TelegramUserChatID int64              `json:"telegram_user_chat_id"`
}

type UserPassword struct {
	ID       primitive.ObjectID `json:"id"`
	Password string             `json:"password"`
}

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserToken struct {
	UserId primitive.ObjectID `json:"id"`
	Token  string             `json:"token"`
}
