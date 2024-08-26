package service

import (
	appError "authservice/internal/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"strings"
	"time"
)

func HandleTelegramMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, serverURL string) {
	switch message.Text {
	case "/start":
		user, err := users.GetUserByTelegramUsername(message.From.UserName)
		if err != nil {
			log.Printf("Ошибка получения userID для Telegram username: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш аккаунт не привязан.")
			bot.Send(msg)
			return
		}

		err = BindTelegramChatID(user.ID, message.Chat.ID)
		if err != nil {
			log.Printf("Ошибка привязки Telegram chat ID: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка привязки Telegram chat ID.")
			bot.Send(msg)
		} else {
			log.Printf("Telegram chat ID %d привязан к userID %s", message.Chat.ID, user.ID.Hex())
			msg := tgbotapi.NewMessage(message.Chat.ID, "Добро пожаловать! Ваш Telegram успешно привязан.")
			bot.Send(msg)
		}

	default:
		user, err := users.GetUserByTelegramUsername(message.From.UserName)
		if err != nil || user.TelegramChatID == 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш аккаунт не привязан.")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Ваш аккаунт уже привязан.")
		bot.Send(msg)
	}
}

func SendCodeToTelegram(chatID int64, code string) error {
	log.Printf("Попытка отправки кода в Telegram для chat ID: %d", chatID)

	msg := tgbotapi.NewMessage(chatID, "Ваш код для входа: "+code)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения в Telegram для chat ID: %d, ошибка: %v", chatID, err)
		if strings.Contains(err.Error(), "Bad Request: chat_id is empty") {
			log.Printf("Ошибка: chat_id пуст. Пользователю необходимо перейти в бота и нажать /start.")

			return appError.ErrChatIDEmpty
		}
		return err
	}

	log.Printf("Сообщение с кодом отправлено в Telegram для chat ID: %d", chatID)
	return nil
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789"
	code := make([]byte, 4)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
