package httphandler

import "go.mongodb.org/mongo-driver/bson/primitive"

type SetUserInfoReq struct {
	Name string `json:"name"`
}

type ChangePswReq struct {
	Password string `json:"password"`
}

type BindTelegramData struct {
	UserID           primitive.ObjectID `json:"user_id"`
	TelegramUsername string             `json:"telegram_username"`
}

type ConfirmTelegramCodeData struct {
	TelegramUsername string `json:"telegram_username"`
	Code             string `json:"code"`
}

type LoginWithTelegramData struct {
	TelegramUsername string `json:"telegram_username"`
}

func (r SetUserInfoReq) IsValid() bool {
	return r.Name != ""
}

func (r ChangePswReq) IsValid() bool {
	return r.Password != ""
}
