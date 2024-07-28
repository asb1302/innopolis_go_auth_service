package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken   string
	ServerURL          string
	CodeExpiryDuration int64 // Время жизни кода в секундах
}

var config Config

func InitConfig() {
	viper.SetDefault("CODE_EXPIRY_DURATION", 300) // По умолчанию 5 минут (300 секунд)

	viper.AutomaticEnv()
	config.TelegramBotToken = viper.GetString("TELEGRAM_BOT_TOKEN")
	config.ServerURL = viper.GetString("SERVER_URL")
	config.CodeExpiryDuration = viper.GetInt64("CODE_EXPIRY_DURATION")
}

func GetConfig() *Config {
	return &config
}
