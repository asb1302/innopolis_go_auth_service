package error

import "errors"

var (
	ErrChatIDEmpty = errors.New("chat_id пуст. Пожалуйста, перейдите в бота и нажмите /start")
)
